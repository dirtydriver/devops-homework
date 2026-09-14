variable "namespace" {
  type    = string
  default = "devops-homework"

  validation {
    condition     = length(trimspace(var.namespace)) > 0
    error_message = "Namespace cannot be empty"
  }
}

variable "environment" {
  description = "Application environment to ENVIRONMENT variable"
  type        = string
  default     = "devops-homework"

  validation {
    condition     = length(trimspace(var.environment)) > 0
    error_message = "Environment cannot be empty"
  }

}

variable "image_repository" {
  description = "Repository of the image"
  type        = string
}
variable "image_tag" {
  description = "Image tag to deploy"
  type        = string
  default     = "1.0.0"

  validation {
    condition     = var.image_tag != "latest"
    error_message = "Cannot use latest tag, you should use fix versions."
  }
}

variable "kubeconfig_path" {
  description = "Path of kubernetes config."
  type        = string
  default     = "~/.kube/config"

}