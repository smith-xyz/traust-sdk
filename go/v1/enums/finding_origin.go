// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// FindingOrigin — Values of `origin` in report.schema.json.
type FindingOrigin string

const (
	FindingOriginVerifyRemediation   FindingOrigin = "verify-remediation"
	FindingOriginVulnScan            FindingOrigin = "vuln-scan"
	FindingOriginCreateFuzzing       FindingOrigin = "create-fuzzing"
	FindingOriginValidateFindings    FindingOrigin = "validate-findings"
	FindingOriginValidationDiscovery FindingOrigin = "validation-discovery"
	FindingOriginImpactAnalysis      FindingOrigin = "impact-analysis"
)

// FindingOriginValues returns all valid FindingOrigin values.
func FindingOriginValues() []FindingOrigin {
	return []FindingOrigin{
		FindingOriginVerifyRemediation,
		FindingOriginVulnScan,
		FindingOriginCreateFuzzing,
		FindingOriginValidateFindings,
		FindingOriginValidationDiscovery,
		FindingOriginImpactAnalysis,
	}
}
