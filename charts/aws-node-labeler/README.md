#### Usage
```terraform
resource "helm_release" "aws-node-labeler" {
  name       = "aws-node-labeler"
  repository = "https://raw.githubusercontent.com/bdwyertech/aws-node-labeler/main/charts"
  chart      = "aws-node-labeler"
}
```
