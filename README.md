# AWS EKS Node Labeler

This application is intended to label Kubernetes nodes with AWS metadata.  By default, this will ensure the label `eks.amazonaws.com/capacityType` is correctly set to `SPOT` or `ON_DEMAND` depending on its `InstanceLifecycle` status. Typically you would set this via a curl in userdata, but AWS Bottlerocket TOML does not support doing this dynamically yet.  This can be a problem if you have an ASG with a mix of spot and on-demand instances.

You can optionally configure this to add other attributes to the tag of your choice.

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
```

Any/all fields in `DescribeInstances` output are available.

https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2/types#Instance

For spot instances, you can use the prefix `instance.spot` to access fields available in `SpotInstanceRequest`.  If an instance is not a spot instance or the field is not available, it will not be set at all.
https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2/types#SpotInstanceRequest


#### AWS Permissions
This application requires `ec2:DescribeInstances` and optionally `ec2:DescribeSpotInstanceRequests`

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
