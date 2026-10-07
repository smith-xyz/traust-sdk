// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

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
