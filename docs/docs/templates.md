---
sidebar_position: 3
---

# Custom Templates

devops-reporter supports custom HTML templates via the `-template` flag. The built-in templates live in `cmd/templates/<source>/template.html` and serve as the reference implementation.

## Usage

```bash
argocd app get my-app -o json | devops-reporter \
  -source argocd \
  -template /path/to/my-template.html
```

## Writing Your Own Template

Templates use Go's [`html/template`](https://pkg.go.dev/html/template) syntax. The data available in the template depends on the source.

### ArgoCD template data

| Field | Type | Description |
|---|---|---|
| `.Title` | string | Report title |
| `.GeneratedAt` | string | UTC timestamp |
| `.AppName` | string | Application name |
| `.AppNamespace` | string | Application namespace |
| `.Project` | string | ArgoCD project |
| `.RepoURL` | string | Source repository URL |
| `.Path` | string | Source path |
| `.TargetRevision` | string | Target revision |
| `.Revision` | string | Actual deployed revision |
| `.DestServer` | string | Destination server |
| `.DestNamespace` | string | Destination namespace |
| `.SyncStatus` | string | Sync status (`Synced`, `OutOfSync`) |
| `.HealthStatus` | string | Health status (`Healthy`, `Degraded`, …) |
| `.OperationPhase` | string | Last operation phase |
| `.OperationMsg` | string | Last operation message |
| `.HasIssues` | bool | `true` when not synced or not healthy |
| `.Summary.Total` | int | Total resource count |
| `.Summary.Synced` | int | Synced count |
| `.Summary.OutOfSync` | int | Out-of-sync count |
| `.Summary.Healthy` | int | Healthy count |
| `.Summary.Degraded` | int | Degraded count |
| `.Summary.Missing` | int | Missing count |
| `.Summary.Unknown` | int | Unknown count |
| `.Groups` | []KindGroup | Resources grouped by kind |
| `.SyncResults` | []SyncResultGroup | Sync op results grouped by kind |
| `.ExternalURLs` | []string | External URLs from app summary |
| `.Images` | []string | Container images from app summary |
| `.SourceType` | string | Source type (e.g. `Kustomize`) |

**Template functions (ArgoCD):**

| Function | Description |
|---|---|
| `syncClass .Status` | Returns a CSS class for sync status |
| `healthClass .Status` | Returns a CSS class for health status |
| `opClass .Phase` | Returns a CSS class for operation phase |
| `shortRev .Revision` | Truncates a git SHA to 7 characters |

### Kubeconform template data

| Field | Type | Description |
|---|---|---|
| `.Title` | string | Report title |
| `.GeneratedAt` | string | UTC timestamp |
| `.TotalCount` | int | Total resources inspected |
| `.HasIssues` | bool | `true` when invalid or errors > 0 |
| `.Summary.Valid` | int | Valid count |
| `.Summary.Invalid` | int | Invalid count |
| `.Summary.Errors` | int | Error count |
| `.Summary.Skipped` | int | Skipped count |
| `.Groups` | []KcKindGroup | Resources grouped by kind |

**Template functions (Kubeconform):**

| Function | Description |
|---|---|
| `statusLabel .Status` | Returns a human-readable label for a status string |
| `statusClass .Status` | Returns a CSS class for a status string |

### Tenable WAS template data

| Field | Type | Description |
|---|---|---|
| `.Title` | string | Report title |
| `.GeneratedAt` | string | UTC timestamp |
| `.ScanTarget` | string | Target URL of the scan |
| `.ScanID` | string | Tenable scan UUID |
| `.ScanStatus` | string | Scan status (e.g. `completed`) |
| `.ScanName` | string | Scan configuration name |
| `.Template` | string | Scan template name (e.g. `basic`) |
| `.StartedAt` | string | Scan start time (UTC) |
| `.FinalizedAt` | string | Scan finish time (UTC) |
| `.Duration` | string | Scan duration (e.g. `3m44s`) |
| `.HasIssues` | bool | `true` when critical/high/medium findings > 0 |
| `.Summary.Critical` | int | Critical finding count |
| `.Summary.High` | int | High finding count |
| `.Summary.Medium` | int | Medium finding count |
| `.Summary.Low` | int | Low finding count |
| `.Summary.Info` | int | Info finding count |
| `.Summary.Total` | int | Total finding count |
| `.Groups` | []WasSeverityGroup | Findings grouped by severity (critical → high → medium → low → info) |

Each item in `.Groups` has:
- `.Severity` — severity string (`critical`, `high`, `medium`, `low`, `info`)
- `.Findings` — slice of findings for that severity

Each finding has `.PluginID`, `.Name`, `.Family`, `.Synopsis`, `.Description`, `.Solution`, `.RiskFactor`, `.URI`, `.CVEs`, `.CWE`, `.OWASP`, `.Output`, `.Proof`, `.SeeAlso`.

**Template functions (Tenable WAS):**

| Function | Description |
|---|---|
| `wasRiskClass .Severity` | Returns a CSS class for a severity string (e.g. `risk-critical`) |
| `wasRiskLabel .Severity` | Returns a human-readable label (e.g. `Critical`) |

### SBOM (CycloneDX) template data

| Field | Type | Description |
|---|---|---|
| `.Title` | string | Report title |
| `.GeneratedAt` | string | UTC timestamp when the report was rendered |
| `.BOMFormat` | string | BOM format name (e.g. `CycloneDX`) |
| `.SpecVersion` | string | CycloneDX spec version (e.g. `1.6`) |
| `.SerialNumber` | string | BOM serial number (URN) |
| `.BOMVersion` | int | BOM document version number |
| `.CreatedAt` | string | Timestamp from the BOM metadata (UTC) |
| `.Lifecycle` | string | Lifecycle phase(s) from metadata (e.g. `build`) |
| `.Tool` | string | Generator tool name and version (first tool in metadata) |
| `.MainComponent` | CdxComponent | The primary application described by the SBOM |
| `.MainComponent.Name` | string | Application name |
| `.MainComponent.Group` | string | Application group/organization |
| `.MainComponent.Version` | string | Application version |
| `.MainComponent.Description` | string | Application description |
| `.MainComponent.PURL` | string | Package URL of the application |
| `.MainComponent.Licenses` | []CdxLicenseEntry | License entries for the application |
| `.MainLicense` | string | Resolved license string for the main component |
| `.HasIssues` | bool | `true` when any component lacks license information |
| `.Summary.Total` | int | Total component count |
| `.Summary.Libraries` | int | Count of components with type `library` |
| `.Summary.Applications` | int | Count of components with type `application` |
| `.Summary.Unlicensed` | int | Count of components with no license info |
| `.Summary.UniqueLicenses` | int | Count of distinct license identifiers |
| `.Summary.Ecosystems` | int | Count of distinct package ecosystems |
| `.Groups` | []CdxEcosystemGroup | Components grouped by package ecosystem |

Each item in `.Groups` has:
- `.Ecosystem` — ecosystem string (e.g. `npm`, `pypi`, `maven`)
- `.Components` — slice of components in that ecosystem, sorted by name

Each component has `.Name`, `.Group`, `.Version`, `.Description`, `.PURL`, `.Type`, `.BOMRef`, `.Licenses`, `.Hashes`.

**Template functions (SBOM CDX):**

| Function | Description |
|---|---|
| `cdxEcosystem .PURL` | Extracts the ecosystem from a PURL (e.g. `npm`, `pypi`) |
| `cdxLicense .Licenses` | Returns the first resolved license string from a `[]CdxLicenseEntry` |
| `cdxShortPurl .PURL` | Strips the `pkg:<ecosystem>/` prefix for compact display |
| `cdxShortHash .Content` | Truncates a hash to the first 12 characters |

### Dependency-Check template data

| Field | Type | Description |
|---|---|---|
| `.Title` | string | Report title |
| `.GeneratedAt` | string | UTC timestamp when the report was rendered |
| `.ProjectName` | string | Project name from `projectInfo.name` |
| `.ReportDate` | string | Report generation date from `projectInfo.reportDate` (UTC) |
| `.EngineVersion` | string | Dependency-Check engine version |
| `.DataSources` | []DepDataSource | NVD / OSS Index data source timestamps |
| `.HasIssues` | bool | `true` when at least one vulnerable dependency exists |
| `.Summary.Total` | int | Total dependency count |
| `.Summary.Vulnerable` | int | Count of dependencies with vulnerabilities |
| `.Summary.Clean` | int | Count of dependencies with no vulnerabilities |
| `.Summary.Critical` | int | Total critical vulnerability count |
| `.Summary.High` | int | Total high vulnerability count |
| `.Summary.Medium` | int | Total medium vulnerability count |
| `.Summary.Low` | int | Total low vulnerability count |
| `.Summary.Info` | int | Total info/negligible vulnerability count |
| `.VulnerableDeps` | []DepDependency | Dependencies with ≥1 vulnerability, sorted by worst severity |
| `.CleanDeps` | []DepDependency | Dependencies with no vulnerabilities |

Each `DepDependency` has `.FileName`, `.FilePath`, `.MD5`, `.SHA1`, `.SHA256`, `.Packages` ([]DepPackage with `.ID`, `.Confidence`, `.URL`), and `.Vulnerabilities`.

Each `DepVulnerability` has `.Source`, `.Name`, `.Severity`, `.CVSSv3` (`*DepCVSSv3` with `.BaseScore`, `.BaseSeverity`, `.AttackVector`), `.CVSSv2` (`*DepCVSSv2` with `.Score`, `.Severity`), `.CWEs`, `.Description`, `.Notes`, `.References` ([]DepReference with `.Source`, `.URL`, `.Name`).

Vulnerabilities within each dependency are pre-sorted by severity (critical first).

**Template functions (Dependency-Check):**

| Function | Description |
|---|---|
| `depSevClass .Severity` | Returns a CSS class for a severity string (e.g. `sev-critical`) |
| `depSevLabel .Severity` | Returns a human-readable label (e.g. `Critical`) |
