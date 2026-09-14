# DevOps Engineer Homework

## Overview

This repository contains a small Go HTTP service packaged as a non-root distroless container and deployable to Kubernetes using Helm and Terraform.

The CI pipeline tests and validates the project, builds the container image, and publishes it to GitHub Container Registry. A manually triggered workflow demonstrates how an existing image could be deployed without rebuilding it.

## Repository Structure

```text
.
├── .github/workflows/
│   ├── ci.yaml
│   └── cd.yaml
├── app/
│   ├── cmd/
│   ├── internal/api/
│   ├── Dockerfile
│   └── go.mod
├── helm/
├── terraform/
├── .gitlab-ci.yml
├── ASSIGNMENT.md
└── REVIEW.md
```

## User Guide

### Prerequisites

* Go 1.27
* Docker
* Helm 3
* kubectl
* Terraform 1.5 or newer
* an existing Kubernetes cluster

### Run locally

```bash
cd app
ENVIRONMENT=local PORT=8080 go run ./cmd
```

Test the service:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/version
curl http://localhost:8080/env
```

### Run tests

```bash
cd app
go vet ./...
go test -race ./...
```

### Build and run the container

```bash
docker build -t devops-homework:1.0.0 ./app

docker run --rm \
  -p 8080:8080 \
  -e ENVIRONMENT=local \
  devops-homework:1.0.0
```

### Deploy with Helm

For a local k3d cluster, import the image before deployment:

```bash
k3d image import devops-homework:1.0.0 \
  --cluster <cluster-name>
```

```bash
helm lint ./helm

helm upgrade --install devops-homework ./helm \
  --namespace devops-homework \
  --create-namespace \
  --set image.repository=devops-homework \
  --set image.tag=1.0.0
```

The service can be tested without an Ingress controller using port forwarding:

```bash
kubectl port-forward \
  -n devops-homework \
  svc/devops-homework-svc \
  8080:80
```

### Deploy with Terraform

```bash
cd terraform

terraform init
terraform validate

terraform apply \
  -var="image_repository=devops-homework" \
  -var="image_tag=1.0.0" \
  -var="environment=local"
```

The Helm and Terraform commands are alternative deployment methods. The same release should not be managed with both at the same time.

## API Endpoints

| Method   | Endpoint         | Description              |
| -------- | ---------------- | ------------------------ |
| `GET`    | `/health`        | Application health       |
| `GET`    | `/version`       | Application version      |
| `GET`    | `/env`           | Current environment      |
| `POST`   | `/config`        | Create or update a value |
| `GET`    | `/config/{name}` | Retrieve a value         |
| `DELETE` | `/config/{name}` | Delete a value           |

Example:

```bash
curl -X POST http://localhost:8080/config \
  -H "Content-Type: application/json" \
  -d '{"name":"database_url","value":"postgres://example"}'
```

## CI/CD

The GitHub Actions CI workflow:

* checks Go formatting and runs `go vet`;
* runs tests with the race detector;
* lints and renders the Helm chart;
* validates the Terraform configuration;
* builds the container image;
* publishes the image to GHCR on pushes to `main`.

Published images use the first 12 characters of the commit SHA as their tag.

The deployment workflow is triggered manually and uses Terraform to deploy an existing image tag. It requires a `KUBECONFIG_BASE64` secret for a Kubernetes cluster reachable from the GitHub runner.

As agreed, I implemented CI/CD in GitHub Actions. The original `.gitlab-ci.yml` remains unchanged.

## What I Changed

* Implemented the API with the Go standard library; the endpoints are simple enough that a framework wasn't needed.
* Kept configuration in memory to avoid a database dependency for local setup, and protected the store with a mutex so concurrent requests can safely access it.
* Added handler tests to check API responses and graceful shutdown to give active requests time to finish when the process stops.
* Used a multi-stage Docker build to keep the compiler out of the runtime image. The service runs as a non-root user in a distroless image.
* Fixed the Helm names, selectors, and ports so the Service routes traffic to the application. Moved deployment settings into `values.yaml` so they can be changed without editing templates.
* Added health probes so Kubernetes can check whether the application is ready and responsive, plus resource requests and limits to set an initial CPU and memory budget.
* Fixed the Terraform providers, variables, and chart path so Terraform can deploy the Helm chart to an existing cluster.
* Added CI checks before publishing images and a manual deployment workflow that reuses a published image without rebuilding it.

## Assumptions

* The Kubernetes cluster already exists.
* A valid kubeconfig is available when deploying.
* An Ingress controller is installed when Ingress is enabled.
* The container registry is accessible from the cluster.
* Values submitted to the configuration API are non-sensitive.
* The resource settings are initial estimates.

## Known Limitations

* Configuration is stored in memory and is lost when the application restarts.
* Configuration is not shared between multiple replicas.
* The configuration API has no authentication.
* No remote Terraform backend is configured.
* The example CD workflow requires an externally provided, reachable cluster.

## Production Improvements

Configuration is lost on restart and isn't shared between pods, so persistent storage would be my first change before scaling the application. The `/config` endpoints also need authentication before exposing the service to other users.

I'd move application deployments to a GitOps workflow with Argo CD, with the deployed image version tracked in Git. That would make application releases easier to manage without running Terraform for each deployment. Terraform would stay responsible for infrastructure, with a remote backend and state locking.

I'd also add request metrics to monitor the service.

## Additional Documents

* [Code Review](REVIEW.md)
* [Original Assignment](ASSIGNMENT.md)
