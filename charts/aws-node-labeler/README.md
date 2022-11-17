#### Usage
```terraform
resource "helm_release" "aws-node-labeler" {
  name       = "aws-node-labeler"
  repository = "oci://ghcr.io/bdwyertech/charts"
  chart      = "bdwyertech"
  version    = "0.1.0"
}
```
