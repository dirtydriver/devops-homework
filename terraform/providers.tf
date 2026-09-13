terraform {
  required_version = ">= 1.5"
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.0"
    }

    helm = {
      source  = "hashicorp/helm"
      version = "~> 3.0"
    }
  }
}




provider "kubernetes" {
  config_path = pathexpand(var.kubeconfig_path)
}

provider "helm" {
  kubernetes = {
    config_path = pathexpand(var.kubeconfig_path)
  }
}