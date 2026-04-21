package main

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"sort"
	"strings"
	"time"
)

//go:embed templates/sbom-cdx.html
var sbomCdxTemplate string

func init() {
	RegisterSource("cyclonedx", &ReportSource{
		DefaultTitle: "Software Bill of Materials",
		Template:     sbomCdxTemplate,
		FuncMap: template.FuncMap{
			"cdxEcosystem": cdxEcosystemFromPURL,
			"cdxLicense":   cdxLicenseString,
			"cdxShortPurl": cdxShortPurl,
			"cdxShortHash": cdxShortHash,
		},
		Parse: func(data []byte, title string) (any, error) {
			var sbom CdxSBOM
			if err := json.Unmarshal(data, &sbom); err != nil {
				return nil, err
			}
			return BuildCdxReportData(sbom, title), nil
		},
	})
}

// ── Input types ──────────────────────────────────────────────────

type CdxLicenseEntry struct {
	License struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"license"`
	Expression string `json:"expression"`
}

type CdxHash struct {
	Alg     string `json:"alg"`
	Content string `json:"content"`
}

type CdxComponent struct {
	Group       string            `json:"group"`
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	PURL        string            `json:"purl"`
	Type        string            `json:"type"`
	BOMRef      string            `json:"bom-ref"`
	Licenses    []CdxLicenseEntry `json:"licenses"`
	Hashes      []CdxHash         `json:"hashes"`
}

type CdxTool struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Group   string `json:"group"`
}

type CdxSBOM struct {
	BOMFormat    string `json:"bomFormat"`
	SpecVersion  string `json:"specVersion"`
	SerialNumber string `json:"serialNumber"`
	Version      int    `json:"version"`
	Metadata     struct {
		Timestamp string `json:"timestamp"`
		Tools     struct {
			Components []CdxTool `json:"components"`
		} `json:"tools"`
		Component  CdxComponent `json:"component"`
		Lifecycles []struct {
			Phase string `json:"phase"`
		} `json:"lifecycles"`
	} `json:"metadata"`
	Components []CdxComponent `json:"components"`
}

// ── Report types ─────────────────────────────────────────────────

type CdxEcosystemGroup struct {
	Ecosystem  string
	Components []CdxComponent
}

type CdxSummary struct {
	Total          int
	Libraries      int
	Applications   int
	Unlicensed     int
	UniqueLicenses int
	Ecosystems     int
}

type CdxReportData struct {
	Title         string
	GeneratedAt   string
	BOMFormat     string
	SpecVersion   string
	SerialNumber  string
	BOMVersion    int
	CreatedAt     string
	Lifecycle     string
	Tool          string
	MainComponent CdxComponent
	MainLicense   string
	Summary       CdxSummary
	Groups        []CdxEcosystemGroup
	HasIssues     bool
}

var cdxEcosystemOrder = []string{"npm", "pypi", "maven", "gem", "cargo", "golang", "nuget", "composer", "other"}

func BuildCdxReportData(sbom CdxSBOM, title string) CdxReportData {
	groupMap := make(map[string][]CdxComponent)
	licenseSet := make(map[string]bool)

	summary := CdxSummary{Total: len(sbom.Components)}

	for _, c := range sbom.Components {
		eco := cdxEcosystemFromPURL(c.PURL)
		groupMap[eco] = append(groupMap[eco], c)

		switch c.Type {
		case "library":
			summary.Libraries++
		case "application":
			summary.Applications++
		}

		lic := cdxLicenseString(c.Licenses)
		if lic == "" {
			summary.Unlicensed++
		} else {
			licenseSet[lic] = true
		}
	}

	summary.UniqueLicenses = len(licenseSet)
	summary.Ecosystems = len(groupMap)

	groups := make([]CdxEcosystemGroup, 0, len(groupMap))
	for _, eco := range cdxEcosystemOrder {
		if comps, ok := groupMap[eco]; ok {
			sort.Slice(comps, func(i, j int) bool {
				return comps[i].Name < comps[j].Name
			})
			groups = append(groups, CdxEcosystemGroup{Ecosystem: eco, Components: comps})
		}
	}

	// any ecosystems not in the canonical order go last (under "other")
	for eco, comps := range groupMap {
		known := false
		for _, k := range cdxEcosystemOrder {
			if k == eco {
				known = true
				break
			}
		}
		if !known {
			sort.Slice(comps, func(i, j int) bool { return comps[i].Name < comps[j].Name })
			groups = append(groups, CdxEcosystemGroup{Ecosystem: eco, Components: comps})
		}
	}

	createdAt := cdxFormatTimestamp(sbom.Metadata.Timestamp)

	tool := ""
	if len(sbom.Metadata.Tools.Components) > 0 {
		t := sbom.Metadata.Tools.Components[0]
		if t.Version != "" {
			tool = t.Name + " " + t.Version
		} else {
			tool = t.Name
		}
	}

	lifecycles := make([]string, 0, len(sbom.Metadata.Lifecycles))
	for _, l := range sbom.Metadata.Lifecycles {
		if l.Phase != "" {
			lifecycles = append(lifecycles, l.Phase)
		}
	}

	return CdxReportData{
		Title:         title,
		GeneratedAt:   time.Now().UTC().Format("2006-01-02 15:04:05 UTC"),
		BOMFormat:     sbom.BOMFormat,
		SpecVersion:   sbom.SpecVersion,
		SerialNumber:  sbom.SerialNumber,
		BOMVersion:    sbom.Version,
		CreatedAt:     createdAt,
		Lifecycle:     strings.Join(lifecycles, ", "),
		Tool:          tool,
		MainComponent: sbom.Metadata.Component,
		MainLicense:   cdxLicenseString(sbom.Metadata.Component.Licenses),
		Summary:       summary,
		Groups:        groups,
		HasIssues:     summary.Unlicensed > 0,
	}
}

func cdxFormatTimestamp(s string) string {
	layouts := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.999999Z",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05.999999-07:00",
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t.UTC().Format("2006-01-02 15:04:05 UTC")
		}
	}
	return s
}

func cdxEcosystemFromPURL(purl string) string {
	if !strings.HasPrefix(purl, "pkg:") {
		return "other"
	}
	rest := purl[4:]
	slash := strings.IndexByte(rest, '/')
	if slash < 0 {
		return "other"
	}
	return rest[:slash]
}

func cdxLicenseString(licenses []CdxLicenseEntry) string {
	for _, l := range licenses {
		if l.License.ID != "" {
			return l.License.ID
		}
		if l.License.Name != "" {
			return l.License.Name
		}
		if l.Expression != "" {
			return l.Expression
		}
	}
	return ""
}

func cdxShortPurl(purl string) string {
	if !strings.HasPrefix(purl, "pkg:") {
		return purl
	}
	rest := purl[4:]
	slash := strings.IndexByte(rest, '/')
	if slash < 0 {
		return purl
	}
	return rest[slash+1:]
}

func cdxShortHash(h string) string {
	if len(h) <= 12 {
		return h
	}
	return h[:12] + "…"
}
