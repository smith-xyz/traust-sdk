// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// ComplianceControlCatalog — Control catalog a control-mapping entry belongs to. FedRAMP is absent on purpose: FedRAMP High and Moderate are 800-53B baseline selections over the nist-800-53-rev5 catalog, so they are assessment targets (compliance_framework), not catalogs.
type ComplianceControlCatalog string

const (
	ComplianceControlCatalogNist80053Rev5 ComplianceControlCatalog = "nist-800-53-rev5"
	ComplianceControlCatalogPciDssV4      ComplianceControlCatalog = "pci-dss-v4"
	ComplianceControlCatalogSoc2Tsc       ComplianceControlCatalog = "soc2-tsc"
	ComplianceControlCatalogGdprTechnical ComplianceControlCatalog = "gdpr-technical"
)

// ComplianceControlCatalogValues returns all valid ComplianceControlCatalog values.
func ComplianceControlCatalogValues() []ComplianceControlCatalog {
	return []ComplianceControlCatalog{
		ComplianceControlCatalogNist80053Rev5,
		ComplianceControlCatalogPciDssV4,
		ComplianceControlCatalogSoc2Tsc,
		ComplianceControlCatalogGdprTechnical,
	}
}
