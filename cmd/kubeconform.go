package main

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"sort"
	"time"
)

//go:embed templates/kubeconform.html
var kubeconformTemplate string

func init() {
	RegisterSource("kubeconform", &ReportSource{
		DefaultTitle: "Kubeconform Validation Report",
		Template:     kubeconformTemplate,
		FuncMap: template.FuncMap{
			"statusLabel": kcStatusLabel,
			"statusClass": kcStatusClass,
		},
		Parse: func(data []byte, title string) (any, error) {
			var output KubeconformOutput
			if err := json.Unmarshal(data, &output); err != nil {
				return nil, err
			}
			return BuildKubeconformReportData(output, title), nil
		},
	})
}

// ── Input types ──────────────────────────────────────────────────

type KcResource struct {
	Filename string `json:"filename"`
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Status   string `json:"status"`
	Msg      string `json:"msg"`
}

type KcSummary struct {
	Valid   int `json:"valid"`
	Invalid int `json:"invalid"`
	Errors  int `json:"errors"`
	Skipped int `json:"skipped"`
}

type KubeconformOutput struct {
	Resources []KcResource `json:"resources"`
	Summary   KcSummary    `json:"summary"`
}

// ── Report types ─────────────────────────────────────────────────

type KcKindGroup struct {
	Kind      string
	Resources []KcResource
}

type KubeconformReportData struct {
	Title       string
	GeneratedAt string
	Summary     KcSummary
	Groups      []KcKindGroup
	TotalCount  int
	HasIssues   bool
}

func BuildKubeconformReportData(output KubeconformOutput, title string) KubeconformReportData {
	kindMap := make(map[string][]KcResource)
	for _, r := range output.Resources {
		kindMap[r.Kind] = append(kindMap[r.Kind], r)
	}

	kinds := make([]string, 0, len(kindMap))
	for k := range kindMap {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)

	groups := make([]KcKindGroup, 0, len(kinds))
	for _, kind := range kinds {
		resources := kindMap[kind]
		sort.Slice(resources, func(i, j int) bool {
			return resources[i].Name < resources[j].Name
		})
		groups = append(groups, KcKindGroup{Kind: kind, Resources: resources})
	}

	return KubeconformReportData{
		Title:       title,
		GeneratedAt: time.Now().UTC().Format("2006-01-02 15:04:05 UTC"),
		Summary:     output.Summary,
		Groups:      groups,
		TotalCount:  len(output.Resources),
		HasIssues:   output.Summary.Invalid+output.Summary.Errors > 0,
	}
}

func kcStatusLabel(status string) string {
	switch status {
	case "statusValid":
		return "Valid"
	case "statusInvalid":
		return "Invalid"
	case "statusError":
		return "Error"
	case "statusSkipped":
		return "Skipped"
	case "statusEmpty":
		return "Empty"
	default:
		return status
	}
}

func kcStatusClass(status string) string {
	switch status {
	case "statusValid":
		return "status-valid"
	case "statusInvalid":
		return "status-invalid"
	case "statusError":
		return "status-error"
	case "statusSkipped":
		return "status-skipped"
	case "statusEmpty":
		return "status-empty"
	default:
		return "status-unknown"
	}
}
