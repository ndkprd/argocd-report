---
sidebar_position: 2
---

# Installation

## Pre-built Binaries

Download the latest release for your platform from the [GitHub Releases](https://github.com/ndkprd/devops-reporter/releases) page.

```bash
# Linux (amd64)
curl -sSL -o devops-reporter \
  https://github.com/ndkprd/devops-reporter/releases/latest/download/devops-reporter_linux_amd64
chmod +x devops-reporter
mv devops-reporter /usr/local/bin/devops-reporter
```

## Build from Source

Requires Go 1.21+.

```bash
git clone https://github.com/ndkprd/devops-reporter.git
cd devops-reporter
go build -o devops-reporter ./cmd/
```

## Docker

```bash
docker build -t devops-reporter .

# Run
docker run --rm -i devops-reporter -source argocd < input.json > report.html
```

Or pull from the GitHub Container Registry:

```bash
docker pull ghcr.io/ndkprd/devops-reporter:latest
```

## Verify

```bash
devops-reporter -version
```
