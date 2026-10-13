// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// ComplianceControlClassification — Values of `classification` in compliance-assessment.schema.json.
type ComplianceControlClassification string

const (
	ComplianceControlClassificationDeterministic  ComplianceControlClassification = "deterministic"
	ComplianceControlClassificationEvidenceReview ComplianceControlClassification = "evidence_review"
	ComplianceControlClassificationOrganizational ComplianceControlClassification = "organizational"
)

// ComplianceControlClassificationValues returns all valid ComplianceControlClassification values.
func ComplianceControlClassificationValues() []ComplianceControlClassification {
	return []ComplianceControlClassification{
		ComplianceControlClassificationDeterministic,
		ComplianceControlClassificationEvidenceReview,
		ComplianceControlClassificationOrganizational,
	}
}
