# Agent Guidelines for devops-reporter

## Project Overview

CLI tool that reads JSON output from DevSecOps tools (ArgoCD, Kubeconform, Tenable WAS, CycloneDX SBOM, OWASP Dependency-Check) and generates self-contained static HTML reports.

## Build Commands

```bash
# Build binary
go build -o devops-reporter ./cmd/

# Build Docker image
docker build -t devops-reporter .

# Run tests (none currently exist)
go test ./...

# Run single test
go test -v -run TestName ./...

# Format code
go fmt ./...

# Vet code
go vet ./...

# Tidy dependencies
go mod tidy
```

## Code Style Guidelines

### Go Version
- Use Go 1.26

### Structure
Each source handler (argocd, kubeconform, tenable-was, sbom-cdx, dep-check) is a separate file under `cmd/` with:
1. `init()` function that registers the source via `RegisterSource()`
2. Input types (JSON binding structs)
3. Report types (internal data structures)
4. Builder function (`Build*ReportData`)
5. Template helper functions

### Type Definitions
```go
// Input types - JSON binding structs use PascalCase fields with json tags
type ArgoResource struct {
    Group     string `json:"group"`
    Version   string `json:"version"`
    Kind      string `json:"kind"`
}

// Report types - internal structures use PascalCase without tags
type ArgoReportData struct {
    Title       string
    Summary     ArgoResourceSummary
    Groups      []ArgoKindGroup
}
```

### Section Headers
Use `// ── Section Name ──` comment separators to organize code:
```go
// ── Input types ──────────────────────────────────────────────────

// ── Report types ─────────────────────────────────────────────────
```

### Error Handling
- Use early returns with `if err != nil`
- Return errors directly from parse functions
- Wrap errors with context where helpful

### Naming Conventions
- Files: lowercase with hyphens (e.g., `sbom-cdx.go`, `tenable-was.go`)
- Types: PascalCase (e.g., `ArgoReportData`, `CdxComponent`)
- Functions: PascalCase for exported, camelCase for unexported
- Template funcs: descriptive names like `syncClass`, `healthClass`, `wasRiskClass`
- Variables: camelCase, prefer meaningful names over abbreviations

### Imports
Standard library only. Use grouping:
```go
import (
    "encoding/json"
    "html/template"
    "sort"
    "time"
)
```

### Common Patterns

**Grouping resources by kind:**
```go
kindMap := make(map[string][]ArgoResource)
for _, r := range app.Status.Resources {
    kindMap[r.Kind] = append(kindMap[r.Kind], r)
}
kinds := make([]string, 0, len(kindMap))
for k := range kindMap {
    kinds = append(kinds, k)
}
sort.Strings(kinds)
groups := make([]ArgoKindGroup, 0, len(kinds))
for _, kind := range kinds {
    resources := kindMap[kind]
    sort.Slice(resources, func(i, j int) bool {
        return resources[i].Name < resources[j].Name
    })
    groups = append(groups, ArgoKindGroup{Kind: kind, Resources: resources})
}
```

**Template helper functions for CSS classes:**
```go
func syncClass(status string) string {
    switch status {
    case "Synced":
        return "sync-synced"
    case "OutOfSync":
        return "sync-outofsync"
    default:
        return "sync-unknown"
    }
}

func healthClass(status string) string {
    switch status {
    case "Healthy":
        return "health-healthy"
    case "Degraded":
        return "health-degraded"
    default:
        return "health-unknown"
    }
}
```

### Templates
- Use `//go:embed` directives to embed HTML templates
- Template files stored in `cmd/templates/` directory
- FuncMap provided for custom template functions

### Docker
- Multi-stage builds
- Builder stage: `golang:1.26-alpine`
- Runtime stage: `alpine:3.23`
- Binary placed at `/usr/local/bin/devops-reporter`

### Testing
- Test files follow `*_test.go` naming
- Use table-driven tests where appropriate
- Test data files stored in `tests/` directory as JSON

## File Structure

```
cmd/
├── main.go              # Entry point, CLI flags, source registration
├── argocd.go            # ArgoCD report handler
├── kubeconform.go       # Kubeconform report handler
├── tenable-was.go       # Tenable WAS report handler
├── sbom-cdx.go          # CycloneDX SBOM report handler
├── dep-check.go         # OWASP Dependency-Check handler
└── templates/           # Embedded HTML templates
    ├── argocd.html
    ├── kubeconform.html
    ├── tenable-was.html
    ├── sbom-cdx.html
    └── dep-check.html
tests/
└── *.json               # Sample input files for each source
```