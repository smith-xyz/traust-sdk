// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// Verdict — Triage finding verdict taxonomy.
type Verdict string

const (
	VerdictTruePositive  Verdict = "true_positive"
	VerdictHardening     Verdict = "hardening"
	VerdictUndetermined  Verdict = "undetermined"
	VerdictFalsePositive Verdict = "false_positive"
	VerdictDuplicate     Verdict = "duplicate"
)

// VerdictValues returns all valid Verdict values.
func VerdictValues() []Verdict {
	return []Verdict{
		VerdictTruePositive,
		VerdictHardening,
		VerdictUndetermined,
		VerdictFalsePositive,
		VerdictDuplicate,
	}
}
