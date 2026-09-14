# Code Review

Review of the skeleton as it was handed over, before my own changes. Top 5 things I'd have this colleague fix before going further.

## 1. Terraform doesn't parse

`terraform/main.tf` has a `set` block with no value:

```hcl
set {
    name  = "image.tag"
    value =
}
```

and further down:

```hcl
set {
    name  = "environment"
    value = prod
}
```

`prod` isn't quoted, so it's read as a reference, not a string — `terraform validate` fails on both. Before touching anything else, Run `terraform fmt`, `terraform init -backend=false`, and `terraform validate` locally or in CI so invalid configuration never reaches a pull request..

## 2. Configuration is inconsistent and partially hard-coded

`namespace = "production"` is hard-coded even though a `namespace` variable already exists and goes unused. The chart path points to a directory that does not exist.

The Helm resources also use inconsistent names and labels: the Deployment uses `app: myapp`, the Service selects `app: myapps`, and the Ingress targets a Service named `homeworks` although the actual Service is named `myapp`. Consequently, traffic cannot reach the application.

Resource names and environment-specific settings should be driven by Terraform variables or Helm values, while selectors and backend references must remain consistent between resources.

## 3. Service routing is broken

The Deployment declares `containerPort: 5000`, while the Service forwards traffic to `targetPort: 8080`. The Ingress correctly targets Service port `80`, but the Service cannot forward that traffic to the declared application port.

The Service selector also uses `app: myapps`, while the pod label is `app: myapp`, so the Service would have no endpoints at all.

Define the application port in `values.yaml`, use it consistently for the container and Service target port, and make sure the Service selector exactly matches the pod labels. Add readiness and liveness probes so Kubernetes can determine whether the application is ready to receive traffic.

## 4. `image.tag: latest` as the default

Use an immutable version such as a Git commit SHA or semantic version and pass it explicitly at deployment time. Avoid `latest` in a chart intended to produce reproducible deployments.

## 5. The pipeline doesn't do anything yet

Both `build` and `deploy` stages in `.gitlab-ci.yml` are `echo` placeholders. There's no lint, no test, no actual image build, no deploy step. Given the Terraform/Helm issues above, I'd get those two working locally first (`terraform validate`, `helm lint` + `helm template`), then wire the same commands into CI so they catch this class of bug automatically next time — right now nothing would have caught points 1-3 before merge.

---

Everything else (structuring the app, Dockerfile, resource requests/limits, etc.) is normal "still needs to be built" territory for a skeleton — these five are the ones that would actively break a deploy or make debugging one painful.
