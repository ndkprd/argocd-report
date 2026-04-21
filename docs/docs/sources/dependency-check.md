---
sidebar_position: 5
---

# OWASP Dependency-Check

The `dependency-check` source reads [OWASP Dependency-Check](https://owasp.org/www-project-dependency-check/) JSON reports and generates a dependency vulnerability analysis report.

## Input

The input must be a Dependency-Check JSON report. Generate it with the `--format JSON` flag:

```bash
# CLI
dependency-check --project my-app --scan . --format JSON --out .

# Maven plugin
mvn dependency-check:check -Dformat=JSON

# Gradle plugin
dependencyCheck { formats = ['JSON'] }
```

## Usage

```bash
cat dependency-check-report.json | devops-reporter -source dependency-check
```

```bash
cat dependency-check-report.json | devops-reporter -source dependency-check \
  -o dep-report.html \
  -title "Dependency Scan — my-app (v1.0.0)"
```

## Report Sections

| Section | Description |
|---|---|
| Header | Report title and project name |
| Scan metadata | Project name, report date, engine version, data source timestamps |
| Status banner | Pass/fail based on presence of vulnerable dependencies |
| Summary grid | Total, Vulnerable, Critical, High, Medium, Low, Info counts |
| Vulnerable dependencies | Each vulnerable dependency as a card showing filename, path, identified packages, and per-CVE details (severity, CVSS score, description, CWEs, references) |
| Clean dependencies | Compact table of all dependencies with no known vulnerabilities |

## HasIssues Flag

`HasIssues` is `true` when at least one dependency has one or more vulnerabilities. This drives the status banner and the header accent color.

## Severity Levels

| Severity | Description |
|---|---|
| Critical | CVSS base score ≥ 9.0 |
| High | CVSS base score 7.0–8.9 |
| Medium | CVSS base score 4.0–6.9 |
| Low | CVSS base score < 4.0 |
| Info | Informational / negligible |

Vulnerabilities within each dependency are sorted by severity (critical first). Vulnerable dependencies are sorted by their worst-severity vulnerability.

## In GitLab CI/CD

```yaml
dependency-check-report:
  stage: test
  image: owasp/dependency-check:latest
  variables:
    DEVOPS_REPORTER_VERSION: v0.2.0
  before_script:
    - |
      wget -qO /usr/local/bin/devops-reporter \
        https://github.com/ndkprd/devops-reporter/releases/download/${DEVOPS_REPORTER_VERSION}/devops-reporter_linux_amd64
      chmod +x /usr/local/bin/devops-reporter
  script:
    - |
      /usr/share/dependency-check/bin/dependency-check.sh \
        --project "${CI_PROJECT_NAME}" \
        --scan . \
        --format JSON \
        --out . \
        --nvdApiKey "${NVD_API_KEY}"
    - |
      cat dependency-check-report.json | devops-reporter -source dependency-check \
        -o dep-report.html \
        -title "Dependency Scan — ${CI_PROJECT_NAME} (${CI_COMMIT_SHORT_SHA})"
  artifacts:
    when: always
    paths:
      - dependency-check-report.json
      - dep-report.html
    expire_in: 1 week
```
