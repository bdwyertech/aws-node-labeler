# AWS EKS Node Labeler

This application is intended to label Kubernetes nodes with AWS metadata.  By default, this will ensure the label `eks.amazonaws.com/capacityType` is correctly set to `SPOT` or `ON_DEMAND` depending on its `InstanceLifecycle` status. Typically you would set this via a curl in userdata, but AWS Bottlerocket TOML does not support doing this dynamically yet.  This can be a problem if you have an ASG with a mix of spot and on-demand instances.

You can optionally configure this to add other attributes to the tag of your choice.

Additionally, this supports appending a suffix to CNI ENI Configuration.  The `aws-vpc-cni` Helm chart [provisions configs per availability zone.](https://github.com/aws/amazon-vpc-cni-k8s/blob/2af69b263885e94e4eeae309b07807b3714c0381/charts/aws-vpc-cni/templates/eniconfig.yaml#L6)  If you have more than one subnet you wish to expose, you can use this feature to dynamically set the correct ENI config.

### Example
```yaml
annotations:
  - name: TeamName
    value: MyCoolTeam
labels:
  - name: aws.bdwyertech.net/zone
    value: instance.Placement.AvailabilityZone
  - name: aws.bdwyertech.net/image
    value: instance.ImageId
  - name: aws.bdwyertech.net/instance
    value: instance.InstanceId
  - name: aws.bdwyertech.net/spotPrice
    value: instance.spot.SpotPrice

eni_config:
  annotation: k8s.amazonaws.com/eniConfig
  suffix_label: aws.bdwyertech.net/eniConfigSuffix
```

#### Result
```
# Annotations
TeamName=MyCoolTeam
# Labels
aws.bdwyertech.net/zone=us-east-1a
aws.bdwyertech.net/image=ami-123456789a9876543
aws.bdwyertech.net/instance=i-abcdef123a456789a
aws.bdwyertech.net/spotPrice=0.768000

# Custom ENI Config (aws.bdwyertech.net/eniConfigSuffix=securedSubnet)
vpc.amazonaws.com/eniConfig=us-east-1c-securedSubnet
```

Any/all fields in `DescribeInstances` output are available.

https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2/types#Instance

For spot instances, you can use the prefix `instance.spot` to access fields available in `SpotInstanceRequest`.  If an instance is not a spot instance or the field is not available, it will not be set at all.
https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2/types#SpotInstanceRequest

A special function, `instance.pod-eni-capable` is available which will set the desired label/annotation to `true` if the instance type is capable of Pod ENI.
https://docs.aws.amazon.com/eks/latest/userguide/security-groups-for-pods.html#supported-instance-types

#### AWS Permissions
##### Required: 
* `ec2:DescribeInstances`

##### Optional:
* `ec2:DescribeSpotInstanceRequests`

#### Kubernetes Permissions
This application requires a ClusterRole similar to the below:
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: aws-node-labeler
  labels:
    app: aws-node-labeler
rules:
  - verbs:
      - get
      - list
      - watch
      - patch
    apiGroups:
      - ""
    resources:
      - nodes
```
