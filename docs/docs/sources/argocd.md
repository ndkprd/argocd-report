---
sidebar_position: 1
---

# ArgoCD

The `argocd` source reads [ArgoCD](https://argoproj.github.io/cd/) Application JSON and generates a deployment status report.

## Input

The input must be ArgoCD Application JSON, as returned by:

```bash
argocd app get <app-name> -o json
# or
kubectl get application <app-name> -n argocd -o json
```

## Usage

```bash
argocd app get my-app -o json | devops-reporter -source argocd
```

```bash
argocd app get my-app -o json | devops-reporter -source argocd \
  -o deploy-report.html \
  -title "Deploy Report — my-app (production)"
```

## Report Sections

| Section | Description |
|---|---|
| Header | App name, namespace, project, source type, resource count |
| Status banner | Sync status and health status at a glance |
| Operation banner | Last sync operation phase and message (if available) |
| Info panel | Source repo/path/revision and destination server/namespace |
| Summary grid | Counts: Synced, Out of Sync, Healthy, Degraded, Missing, Unknown |
| Resources by kind | All resources grouped by Kubernetes Kind, with sync and health badges |
| Sync operation results | Per-resource results of the last sync operation |
| Artifacts | External URLs and container images from the app summary |

## In GitLab CI/CD

```yaml
generate-deploy-report:
  stage: report
  image: quay.io/argoproj/argocd:latest
  variables:
    DEVOPS_REPORTER_VERSION: v0.2.0
  before_script:
    - |
      curl -sSL -o /usr/local/bin/devops-reporter \
        https://github.com/ndkprd/devops-reporter/releases/download/${DEVOPS_REPORTER_VERSION}/devops-reporter_linux_amd64
      chmod +x /usr/local/bin/devops-reporter
  script:
    - |
      argocd app get ${ARGOCD_APP_NAME} -o json | devops-reporter \
        -source argocd \
        -o report.html \
        -title "Deploy Report — ${CI_PROJECT_NAME} (${CI_ENVIRONMENT_NAME})"
  artifacts:
    when: always
    paths:
      - report.html
    expire_in: 1 week
```
