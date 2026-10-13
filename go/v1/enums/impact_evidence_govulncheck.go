// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// ImpactEvidenceGovulncheck — Govulncheck reachability classification.
type ImpactEvidenceGovulncheck string

const (
	ImpactEvidenceGovulncheckSymbolReachable            ImpactEvidenceGovulncheck = "symbol_reachable"
	ImpactEvidenceGovulncheckPackageImportedNotObserved ImpactEvidenceGovulncheck = "package_imported_not_observed"
	ImpactEvidenceGovulncheckModuleRequiredNotObserved  ImpactEvidenceGovulncheck = "module_required_not_observed"
)

// ImpactEvidenceGovulncheckValues returns all valid ImpactEvidenceGovulncheck values.
func ImpactEvidenceGovulncheckValues() []ImpactEvidenceGovulncheck {
	return []ImpactEvidenceGovulncheck{
		ImpactEvidenceGovulncheckSymbolReachable,
		ImpactEvidenceGovulncheckPackageImportedNotObserved,
		ImpactEvidenceGovulncheckModuleRequiredNotObserved,
	}
}
