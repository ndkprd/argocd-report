---
sidebar_position: 2
---

# Kubeconform

The `kubeconform` source reads [kubeconform](https://github.com/yannh/kubeconform) JSON output and generates a Kubernetes manifest validation report.

## Input

The input must be kubeconform JSON output, produced with the `-output json` flag:

```bash
kubeconform -output json ./manifests/
```

## Usage

```bash
kubeconform -output json ./manifests/ | devops-reporter -source kubeconform
```

```bash
kubeconform -output json ./manifests/ | devops-reporter -source kubeconform \
  -o validation-report.html \
  -title "Manifest Validation — my-app (staging)"
```

## Report Sections

| Section | Description |
|---|---|
| Header | Report title and total resources inspected |
| Status banner | Pass/fail overall result |
| Summary grid | Counts: Valid, Invalid, Errors, Skipped |
| Resources by kind | All resources grouped by Kubernetes Kind, with per-resource status badges and messages |

## Status Values

| Status | Meaning |
|---|---|
| Valid | Resource passed schema validation |
| Invalid | Resource failed schema validation |
| Error | Could not validate (e.g. schema not found for CRD) |
| Skipped | Resource was skipped |
| Empty | File contained no resources |

## In GitLab CI/CD

```yaml
validate-manifests:
  stage: validate
  image: ghcr.io/yannh/kubeconform:latest
  variables:
    DEVOPS_REPORTER_VERSION: v0.2.0
  before_script:
    - |
      curl -sSL -o /usr/local/bin/devops-reporter \
        https://github.com/ndkprd/devops-reporter/releases/download/${DEVOPS_REPORTER_VERSION}/devops-reporter_linux_amd64
      chmod +x /usr/local/bin/devops-reporter
  script:
    - |
      kubeconform -output json ./manifests/ | devops-reporter \
        -source kubeconform \
        -o validation-report.html \
        -title "Manifest Validation — ${CI_PROJECT_NAME}"
  artifacts:
    when: always
    paths:
      - validation-report.html
    expire_in: 1 week
```
