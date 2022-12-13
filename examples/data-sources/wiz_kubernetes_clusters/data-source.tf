# Get Azure Cloud hosted Kubernetes Clusters 
data "wiz_kubernetes_clusters" "myclusters" {
  kind = ["AKS"]
}

# Get the first 3 clusters on a specific AWS account ID
data "wiz_kubernetes_clusters" "myclusters" {
  external_ids = ["232412319201"]
  first        = 3
}