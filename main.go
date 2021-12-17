// Encoding: UTF-8
//
// AWS Node Labeler
//
// Copyright © 2021 Brian Dwyer - Intelligent Digital Services
//

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Jeffail/gabs/v2"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"gopkg.in/yaml.v3"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	// Kubernetes
	"github.com/aws/amazon-vpc-resource-controller-k8s/pkg/aws/vpc"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

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

type mutator struct {
	client *kubernetes.Clientset
	config Config
	ctx    context.Context
}

type Config struct {
	Annotations []keyValue `yaml:"annotations"`
	Labels      []keyValue `yaml:"labels"`
	EniConfig   *struct {
		Label       string `yaml:"label"`
		SuffixLabel string `yaml:"suffix_label"`
	} `yaml:"eni_config"`
}

type keyValue struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// https://firehydrant.io/blog/stay-informed-with-kubernetes-informers/

func main() {
	flag.Parse()

	if versionFlag {
		showVersion()
		os.Exit(0)
	}
	var cfg Config
	if val, ok := os.LookupEnv("CONFIG_FILE"); ok {
		cfgFile, err := os.Open(val)
		if err != nil {
			log.Fatal(err)
		}

		cfgBytes, err := io.ReadAll(cfgFile)
		if err != nil {
			log.Fatal(err)
		}
		cfgFile.Close()

		if err = yaml.Unmarshal(cfgBytes, &cfg); err != nil {
			log.Fatal(err)
		}
	}

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		log.Fatal(err)
	}
	client, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("Target:", kubeconfig.Host)

	mu := &mutator{client, cfg, context.Background()}

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
	log *log.Entry
}

func (n *Node) Annotate(key, value string) {
	annotations := n.GetAnnotations()
	if val, ok := annotations[key]; !ok || val != value {
		n.log.Infof("Setting Annotation: %s=%s", key, value)
		annotations[key] = value
	}
}

func (n *Node) Label(key, value string) {
	value = regexp.MustCompile("[^a-zA-Z0-9-_.]+").ReplaceAllString(value, "-")
	matches := regexp.MustCompile(`(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?`).FindAllString(value, -1)
	for _, value = range matches {
		if value != "" {
			break
		}
	}
	labels := n.GetLabels()
	if val, ok := labels[key]; !ok || val != value {
		n.log.Infof("Setting Label: %s=%s", key, value)
		labels[key] = value
	}
}

func (mu *mutator) Add(obj interface{}) {
	nodeObj := obj.(*corev1.Node)
	nodeName := nodeObj.GetName()
	log := log.WithFields(log.Fields{"node": nodeName})
	node := Node{nodeObj, log}

	if val, ok := node.GetLabels()["eks.amazonaws.com/compute-type"]; ok {
		if val == "fargate" {
			log.Debugln("Skipping fargate node:", nodeName)
			return
		}
	}

	oldData, err := json.Marshal(nodeObj)
	if err != nil {
		log.Error(err)
		return
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

	if eniConfig := mu.config.EniConfig; eniConfig != nil {
		if val, ok := node.GetLabels()[eniConfig.SuffixLabel]; ok {
			node.Label(eniConfig.Label, fmt.Sprintf("%s-%s", *instance.Placement.AvailabilityZone, val))
		}
	}

	jsonBytes, err := json.Marshal(instance)
	if err != nil {
		log.Fatal(err)
	}
	instanceObj, err := gabs.ParseJSON(jsonBytes)
	if err != nil {
		log.Fatal(err)
	}

	var spotOnce sync.Once
	var spotObj *gabs.Container
	getSpot := func() *gabs.Container {
		spotOnce.Do(func() {
			// Get Spot Instance Details
			if reservations, err := ec2Client.DescribeSpotInstanceRequests(mu.ctx, &ec2.DescribeSpotInstanceRequestsInput{
				SpotInstanceRequestIds: []string{*instance.SpotInstanceRequestId},
			}); err != nil {
				log.Error(err)
				return
			} else if len(reservations.SpotInstanceRequests) == 1 {
				jsonBytes, err := json.Marshal(reservations.SpotInstanceRequests[0])
				if err != nil {
					log.Error(err)
					return
				}
				spotObj, err = gabs.ParseJSON(jsonBytes)
				if err != nil {
					log.Error(err)
				}
			}
		})
		return spotObj
	}

	apply := func(applyFunc func(string, string), kv []keyValue) {
		for _, v := range kv {
			if v.Value == "instance.pod-eni-capable" {
				// Pod ENI is supported by many different instance types
				// https://docs.aws.amazon.com/eks/latest/userguide/security-groups-for-pods.html#supported-instance-types
				// VPC Resource Controller - IsInstanceSupported
				// https://github.com/aws/amazon-vpc-resource-controller-k8s/blob/05e89ff9300a5cc0ebea705834cf27f0a7f3b509/pkg/provider/branch/provider.go#L463-L472
				limits, found := vpc.Limits[string(instance.InstanceType)]
				if found && instance.Platform != ec2Types.PlatformValuesWindows && limits.IsTrunkingCompatible {
					applyFunc(v.Name, "true")
				}
			} else if pfx := "instance.spot."; strings.HasPrefix(v.Value, pfx) {
				if instance.SpotInstanceRequestId == nil {
					continue
				}
				if spotObj := getSpot(); spotObj != nil {
					if val := spotObj.Path(strings.TrimPrefix(v.Value, pfx)).Data(); val != nil {
						applyFunc(v.Name, val.(string))
					}
				}
			} else if pfx := "instance."; strings.HasPrefix(v.Value, pfx) {
				if val := instanceObj.Path(strings.TrimPrefix(v.Value, pfx)).Data(); val != nil {
					applyFunc(v.Name, val.(string))
				}
			} else {
				applyFunc(v.Name, v.Value)
			}
		}
	}
	apply(node.Annotate, mu.config.Annotations)
	apply(node.Label, mu.config.Labels)

	newData, err := json.Marshal(nodeObj)
	if err != nil {
		log.Error(err)
		return
	}

	if !reflect.DeepEqual(oldData, newData) {
		patchBytes, err := jsonpatch.CreateMergePatch(oldData, newData)
		if err != nil {
			log.Error(err)
			return
		}
		_, err = mu.client.CoreV1().Nodes().Patch(mu.ctx, nodeName, k8stypes.MergePatchType, patchBytes, metav1.PatchOptions{})
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("Updated Node")
	}
}
