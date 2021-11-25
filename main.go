package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	// Kubernetes
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type mutator struct {
	client *kubernetes.Clientset
	ctx    context.Context
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	// Debug Logging
	if _, ok := os.LookupEnv("DEBUG"); ok {
		log.SetLevel(log.DebugLevel)
	}
	if _, ok := os.LookupEnv("DEBUG_TRACE"); ok {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}
}

// https://firehydrant.io/blog/stay-informed-with-kubernetes-informers/

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		log.Fatal(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("Target:", config.Host)

	mu := &mutator{client, context.Background()}

	factory := informers.NewSharedInformerFactory(client, 0)
	informer := factory.Core().V1().Nodes().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: mu.Add,
	})
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

type Node struct {
	*corev1.Node
	updated bool
	log     *log.Entry
}

func (n *Node) Label(key, value string) {
	labels := n.GetLabels()
	if val, ok := labels[key]; !ok || val != value {
		n.log.Infof("Setting Label: %s=%s", key, value)
		labels[key] = value
		n.updated = true
	}
}

func (mu *mutator) Add(obj interface{}) {
	nodeObj := obj.(*corev1.Node)
	nodeName := nodeObj.GetName()
	log := log.WithFields(log.Fields{"node": nodeName})
	node := Node{nodeObj, false, log}

	if val, ok := node.GetLabels()["eks.amazonaws.com/compute-type"]; ok {
		if val == "fargate" {
			log.Debugln("Skipping fargate node:", nodeName)
			return
		}
	}

	// ProviderID
	// EC2 - aws:///us-east-1c/i-0e190165ce4facc0f
	// Fargate - aws:///us-east-1b/b7af340c11-9ec3eeb6643c4ea58b0285cefd83ef94/fargate-ip-10-65-48-87.ec2.internal
	if !strings.HasPrefix(node.Spec.ProviderID, "aws:/") {
		log.Debug("Not an AWS Node... Skipping.")
		return
	}
	instanceID := filepath.Base(node.Spec.ProviderID)
	instanceAz := strings.Split(strings.TrimPrefix(node.Spec.ProviderID, "aws:///"), "/")[0]
	var region string
	if l := len(instanceAz); l > 0 {
		region = instanceAz[:l-1]
	}
	// region, ok := node.GetLabels()["failure-domain.beta.kubernetes.io/region"]
	// if !ok {
	// 	log.Error("Region not found")
	// 	return
	// }

	awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Error(err)
		return
	}

	ec2Client := ec2.NewFromConfig(awsCfg)

	instancesOutput, err := ec2Client.DescribeInstances(mu.ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		log.Error(err)
		return
	}
	if len(instancesOutput.Reservations) != 1 {
		log.Error("expected one reservation")
		return
	}
	if len(instancesOutput.Reservations[0].Instances) != 1 {
		log.Error("expected one instance")
		return
	}
	instance := instancesOutput.Reservations[0].Instances[0]

	lifecycle := "ON_DEMAND"
	switch strings.ToUpper(string(instance.InstanceLifecycle)) {
	case "SPOT":
		lifecycle = "SPOT"
	}

	node.Label("eks.amazonaws.com/capacityType", lifecycle)

	if node.updated {
		nodeObj.SetLabels(node.GetLabels())
		_, err = mu.client.CoreV1().Nodes().Update(mu.ctx, nodeObj, metav1.UpdateOptions{})
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("Updated Node")
	}
}
