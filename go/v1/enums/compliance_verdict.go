// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// ComplianceVerdict — Compliance control assessment verdict.
type ComplianceVerdict string

const (
	ComplianceVerdictSatisfied     ComplianceVerdict = "satisfied"
	ComplianceVerdictNotSatisfied  ComplianceVerdict = "not_satisfied"
	ComplianceVerdictNotApplicable ComplianceVerdict = "not_applicable"
	ComplianceVerdictNotAssessed   ComplianceVerdict = "not_assessed"
)

// ComplianceVerdictValues returns all valid ComplianceVerdict values.
func ComplianceVerdictValues() []ComplianceVerdict {
	return []ComplianceVerdict{
		ComplianceVerdictSatisfied,
		ComplianceVerdictNotSatisfied,
		ComplianceVerdictNotApplicable,
		ComplianceVerdictNotAssessed,
	}
}
