---
sidebar_position: 7
---

# Trivy

The `trivy` source reads [Trivy](https://trivy.dev/) container image scan JSON output (schema version 2) and generates a vulnerability and package inventory report.

## Input

The input must be a Trivy JSON report produced with `-f json`. Run a full image scan:

```bash
trivy image -f json -o trivy-report.json my-image:tag
```

Or pipe directly:

```bash
trivy image -f json my-image:tag | devops-reporter -source trivy
```

## Usage

```bash
cat trivy-report.json | devops-reporter -source trivy
```

```bash
cat trivy-report.json | devops-reporter -source trivy \
  -o trivy-report.html \
  -title "Vulnerability Scan — my-app (v1.0.0)"
```

## Report Sections

| Section | Description |
|---|---|
| Header | Report title and artifact name |
| Scan metadata | Artifact type, OS, scan date, Trivy version, target count, fixable ratio |
| Image details | Size, architecture, build timestamp, layer count, image ID, revision, maintainer, source URL |
| Status banner | Pass/fail based on presence of any vulnerabilities |
| Summary grid | Total, Critical, High, Medium, Low, Fixable counts |
| Vulnerabilities by severity | Per-severity groups with CVE ID (linked), package, installed version, fixed version, status, description |
| Installed packages | Per-target groups listing all packages with name, version, license(s), source package, and layer digest |
| Footer | Full count summary including package total |

## HasIssues Flag

`HasIssues` is `true` when the report contains at least one vulnerability. This drives the status banner and the header accent color.

## Severity Levels

| Severity | Description |
|---|---|
| Critical | Immediate risk — exploit exists or CVSS ≥ 9.0 |
| High | Significant impact — CVSS 7.0–8.9 |
| Medium | Moderate impact — CVSS 4.0–6.9 |
| Low | Minimal direct impact — CVSS < 4.0 |
| Unknown | Severity not yet assessed |

## In GitLab CI/CD

```yaml
trivy-scan:
  stage: test
  image:
    name: aquasec/trivy:latest
    entrypoint: [""]
  variables:
    DEVOPS_REPORTER_VERSION: v0.2.0
    IMAGE: $CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA
  before_script:
    - |
      wget -qO /usr/local/bin/devops-reporter \
        https://github.com/ndkprd/devops-reporter/releases/download/${DEVOPS_REPORTER_VERSION}/devops-reporter_linux_amd64
      chmod +x /usr/local/bin/devops-reporter
  script:
    - |
      trivy image -f json -o trivy-report.json ${IMAGE}
    - |
      cat trivy-report.json | devops-reporter -source trivy \
        -o trivy-report.html \
        -title "Vulnerability Scan — ${CI_PROJECT_NAME} (${CI_COMMIT_REF_NAME})"
  artifacts:
    when: always
    paths:
      - trivy-report.json
      - trivy-report.html
    expire_in: 1 week
```
