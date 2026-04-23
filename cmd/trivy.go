package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"sort"
	"strings"
	"time"
)

//go:embed templates/trivy.html
var trivyTemplate string

func init() {
	RegisterSource("trivy", &ReportSource{
		DefaultTitle: "Trivy Vulnerability Report",
		Template:     trivyTemplate,
		FuncMap: template.FuncMap{
			"trivySevClass":    trivySevClass,
			"trivySevLabel":    trivySevLabel,
			"trivyJoinLicenses": trivyJoinLicenses,
			"trivyShortHash":   trivyShortHash,
		},
		Parse: func(data []byte, title string) (any, error) {
			var report TrivyReport
			if err := json.Unmarshal(data, &report); err != nil {
				return nil, err
			}
			return BuildTrivyReportData(report, title), nil
		},
	})
}

// ── Input types ──────────────────────────────────────────────────

type TrivyVulnerability struct {
	VulnerabilityID  string   `json:"VulnerabilityID"`
	PkgID            string   `json:"PkgID"`
	PkgName          string   `json:"PkgName"`
	InstalledVersion string   `json:"InstalledVersion"`
	FixedVersion     string   `json:"FixedVersion"`
	Status           string   `json:"Status"`
	PrimaryURL       string   `json:"PrimaryURL"`
	Description      string   `json:"Description"`
	Severity         string   `json:"Severity"`
	References       []string `json:"References"`
}

type TrivyPackage struct {
	ID         string   `json:"ID"`
	Name       string   `json:"Name"`
	Version    string   `json:"Version"`
	Arch       string   `json:"Arch"`
	SrcName    string   `json:"SrcName"`
	SrcVersion string   `json:"SrcVersion"`
	Licenses   []string `json:"Licenses"`
	Maintainer string   `json:"Maintainer"`
	Digest     string   `json:"Digest"`
	Layer      struct {
		DiffID string `json:"DiffID"`
	} `json:"Layer"`
}

type TrivyLayer struct {
	Size   int64  `json:"Size"`
	DiffID string `json:"DiffID"`
}

type TrivyResult struct {
	Target          string               `json:"Target"`
	Class           string               `json:"Class"`
	Type            string               `json:"Type"`
	Packages        []TrivyPackage       `json:"Packages"`
	Vulnerabilities []TrivyVulnerability `json:"Vulnerabilities"`
}

type TrivyReport struct {
	SchemaVersion int    `json:"SchemaVersion"`
	ReportID      string `json:"ReportID"`
	CreatedAt     string `json:"CreatedAt"`
	ArtifactID    string `json:"ArtifactID"`
	ArtifactName  string `json:"ArtifactName"`
	ArtifactType  string `json:"ArtifactType"`
	Trivy         struct {
		Version string `json:"Version"`
	} `json:"Trivy"`
	Metadata struct {
		Size int64 `json:"Size"`
		OS   struct {
			Family string `json:"Family"`
			Name   string `json:"Name"`
		} `json:"OS"`
		ImageID     string      `json:"ImageID"`
		RepoTags    []string    `json:"RepoTags"`
		RepoDigests []string    `json:"RepoDigests"`
		Layers      []TrivyLayer `json:"Layers"`
		ImageConfig struct {
			Architecture string `json:"architecture"`
			Created      string `json:"created"`
			OS           string `json:"os"`
			Config       struct {
				Labels map[string]string `json:"Labels"`
			} `json:"config"`
		} `json:"ImageConfig"`
	} `json:"Metadata"`
	Results []TrivyResult `json:"Results"`
}

// ── Report types ─────────────────────────────────────────────────

type TrivySeverityGroup struct {
	Severity string
	Vulns    []TrivyVulnerability
}

type TrivyPackageGroup struct {
	Target   string
	Type     string
	Class    string
	Packages []TrivyPackage
}

type TrivyImageDetails struct {
	Size         string
	Architecture string
	BuiltAt      string
	ImageID      string
	RepoDigest   string
	Maintainer   string
	Source       string
	Revision     string
	LayerCount   int
}

type TrivySummary struct {
	Total    int
	Critical int
	High     int
	Medium   int
	Low      int
	Unknown  int
	Fixable  int
	Targets  int
	Packages int
}

type TrivyReportData struct {
	Title          string
	GeneratedAt    string
	ArtifactName   string
	ArtifactType   string
	TrivyVersion   string
	CreatedAt      string
	OS             string
	ImageID        string
	ImageDetails   TrivyImageDetails
	Summary        TrivySummary
	PackageGroups  []TrivyPackageGroup
	SeverityGroups []TrivySeverityGroup
	HasIssues      bool
}

var trivySeverityOrder = []string{"CRITICAL", "HIGH", "MEDIUM", "LOW", "UNKNOWN"}


func BuildTrivyReportData(report TrivyReport, title string) TrivyReportData {
	summary := TrivySummary{Targets: len(report.Results)}
	sevMap := make(map[string][]TrivyVulnerability)
	var packageGroups []TrivyPackageGroup

	for _, result := range report.Results {
		for _, v := range result.Vulnerabilities {
			sev := strings.ToUpper(v.Severity)
			sevMap[sev] = append(sevMap[sev], v)
			summary.Total++

			switch sev {
			case "CRITICAL":
				summary.Critical++
			case "HIGH":
				summary.High++
			case "MEDIUM":
				summary.Medium++
			case "LOW":
				summary.Low++
			default:
				summary.Unknown++
			}

			if v.FixedVersion != "" {
				summary.Fixable++
			}
		}

		if len(result.Packages) > 0 {
			pkgs := make([]TrivyPackage, len(result.Packages))
			copy(pkgs, result.Packages)
			sort.Slice(pkgs, func(i, j int) bool {
				return pkgs[i].Name < pkgs[j].Name
			})
			summary.Packages += len(pkgs)
			packageGroups = append(packageGroups, TrivyPackageGroup{
				Target:   result.Target,
				Type:     result.Type,
				Class:    result.Class,
				Packages: pkgs,
			})
		}
	}

	severityGroups := make([]TrivySeverityGroup, 0, len(sevMap))
	for _, sev := range trivySeverityOrder {
		if vulns, ok := sevMap[sev]; ok {
			sort.Slice(vulns, func(i, j int) bool {
				return vulns[i].PkgName < vulns[j].PkgName
			})
			severityGroups = append(severityGroups, TrivySeverityGroup{Severity: sev, Vulns: vulns})
		}
	}

	os := ""
	if report.Metadata.OS.Family != "" {
		os = report.Metadata.OS.Family
		if report.Metadata.OS.Name != "" {
			os += " " + report.Metadata.OS.Name
		}
	}

	labels := report.Metadata.ImageConfig.Config.Labels
	revision := labels["org.opencontainers.image.revision"]
	if len(revision) > 7 {
		revision = revision[:7]
	}

	repoDigest := ""
	if len(report.Metadata.RepoDigests) > 0 {
		repoDigest = report.Metadata.RepoDigests[0]
	}

	imageDetails := TrivyImageDetails{
		Size:         trivyFormatSize(report.Metadata.Size),
		Architecture: report.Metadata.ImageConfig.Architecture,
		BuiltAt:      trivyFormatTimestamp(report.Metadata.ImageConfig.Created),
		ImageID:      trivyShortHash(report.Metadata.ImageID),
		RepoDigest:   repoDigest,
		Maintainer:   labels["maintainer"],
		Source:       labels["org.opencontainers.image.source"],
		Revision:     revision,
		LayerCount:   len(report.Metadata.Layers),
	}

	createdAt := trivyFormatTimestamp(report.CreatedAt)

	return TrivyReportData{
		Title:          title,
		GeneratedAt:    time.Now().UTC().Format("2006-01-02 15:04:05 UTC"),
		ArtifactName:   report.ArtifactName,
		ArtifactType:   report.ArtifactType,
		TrivyVersion:   report.Trivy.Version,
		CreatedAt:      createdAt,
		OS:             os,
		ImageID:        report.Metadata.ImageID,
		ImageDetails:   imageDetails,
		Summary:        summary,
		PackageGroups:  packageGroups,
		SeverityGroups: severityGroups,
		HasIssues:      summary.Total > 0,
	}
}

func trivyFormatTimestamp(s string) string {
	layouts := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.999999999Z",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05.999999999-07:00",
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t.UTC().Format("2006-01-02 15:04:05 UTC")
		}
	}
	return s
}

func trivyFormatSize(bytes int64) string {
	if bytes == 0 {
		return ""
	}
	const (
		mb = 1024 * 1024
		gb = 1024 * mb
	)
	if bytes >= gb {
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(gb))
	}
	return fmt.Sprintf("%.1f MB", float64(bytes)/float64(mb))
}

func trivyShortHash(h string) string {
	// Strip "sha256:" prefix for display
	h = strings.TrimPrefix(h, "sha256:")
	if len(h) <= 19 {
		return h
	}
	return h[:19] + "…"
}

func trivyJoinLicenses(ls []string) string {
	if len(ls) == 0 {
		return ""
	}
	return strings.Join(ls, ", ")
}

func trivySevClass(severity string) string {
	switch strings.ToUpper(severity) {
	case "CRITICAL":
		return "sev-critical"
	case "HIGH":
		return "sev-high"
	case "MEDIUM":
		return "sev-medium"
	case "LOW":
		return "sev-low"
	default:
		return "sev-unknown"
	}
}

func trivySevLabel(severity string) string {
	switch strings.ToUpper(severity) {
	case "CRITICAL":
		return "Critical"
	case "HIGH":
		return "High"
	case "MEDIUM":
		return "Medium"
	case "LOW":
		return "Low"
	default:
		return "Unknown"
	}
}
