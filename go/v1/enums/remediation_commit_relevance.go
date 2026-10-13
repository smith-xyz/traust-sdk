// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// RemediationCommitRelevance — Values of `relevance` in verification.schema.json.
type RemediationCommitRelevance string

const (
	RemediationCommitRelevanceDirect     RemediationCommitRelevance = "direct"
	RemediationCommitRelevanceSupporting RemediationCommitRelevance = "supporting"
	RemediationCommitRelevancePartial    RemediationCommitRelevance = "partial"
)

// RemediationCommitRelevanceValues returns all valid RemediationCommitRelevance values.
func RemediationCommitRelevanceValues() []RemediationCommitRelevance {
	return []RemediationCommitRelevance{
		RemediationCommitRelevanceDirect,
		RemediationCommitRelevanceSupporting,
		RemediationCommitRelevancePartial,
	}
}
