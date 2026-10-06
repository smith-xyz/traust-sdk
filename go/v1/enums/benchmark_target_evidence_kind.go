// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// BenchmarkTargetEvidenceKind — Values of `kind` in benchmark-target.schema.json.
type BenchmarkTargetEvidenceKind string

const (
	BenchmarkTargetEvidenceKindVerificationReport BenchmarkTargetEvidenceKind = "verification_report"
	BenchmarkTargetEvidenceKindLedgerEvent        BenchmarkTargetEvidenceKind = "ledger_event"
	BenchmarkTargetEvidenceKindAdvisory           BenchmarkTargetEvidenceKind = "advisory"
	BenchmarkTargetEvidenceKindFuzzCrash          BenchmarkTargetEvidenceKind = "fuzz_crash"
	BenchmarkTargetEvidenceKindSeedRecipe         BenchmarkTargetEvidenceKind = "seed_recipe"
)

// BenchmarkTargetEvidenceKindValues returns all valid BenchmarkTargetEvidenceKind values.
func BenchmarkTargetEvidenceKindValues() []BenchmarkTargetEvidenceKind {
	return []BenchmarkTargetEvidenceKind{
		BenchmarkTargetEvidenceKindVerificationReport,
		BenchmarkTargetEvidenceKindLedgerEvent,
		BenchmarkTargetEvidenceKindAdvisory,
		BenchmarkTargetEvidenceKindFuzzCrash,
		BenchmarkTargetEvidenceKindSeedRecipe,
	}
}
