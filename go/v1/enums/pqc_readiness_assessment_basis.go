// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// PqcReadinessAssessmentBasis — Values of `assessment_basis` in pqc-readiness.schema.json.
type PqcReadinessAssessmentBasis string

const (
	PqcReadinessAssessmentBasisSource        PqcReadinessAssessmentBasis = "source"
	PqcReadinessAssessmentBasisSbomOnly      PqcReadinessAssessmentBasis = "sbom-only"
	PqcReadinessAssessmentBasisSourceRuntime PqcReadinessAssessmentBasis = "source+runtime"
)

// PqcReadinessAssessmentBasisValues returns all valid PqcReadinessAssessmentBasis values.
func PqcReadinessAssessmentBasisValues() []PqcReadinessAssessmentBasis {
	return []PqcReadinessAssessmentBasis{
		PqcReadinessAssessmentBasisSource,
		PqcReadinessAssessmentBasisSbomOnly,
		PqcReadinessAssessmentBasisSourceRuntime,
	}
}
