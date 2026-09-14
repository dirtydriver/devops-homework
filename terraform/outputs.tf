output "namespace" {
  description = "The namespace where the application was deployed"
  value       = kubernetes_namespace_v1.homework.metadata[0].name
}
output "helm_release_name" {
  description = "Name of the helm release"
  value       = helm_release.homework.name
}

output "helm_release_status" {
  description = "Status of the Helm Release"
  value       = helm_release.homework.status
}