// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// RestatementReason — Why the data was wrong.
type RestatementReason string

const (
	RestatementReasonDataError       RestatementReason = "data_error"
	RestatementReasonSchemaMigration RestatementReason = "schema_migration"
	RestatementReasonBaselineRewrite RestatementReason = "baseline_rewrite"
	RestatementReasonOperatorError   RestatementReason = "operator_error"
)

// RestatementReasonValues returns all valid RestatementReason values.
func RestatementReasonValues() []RestatementReason {
	return []RestatementReason{
		RestatementReasonDataError,
		RestatementReasonSchemaMigration,
		RestatementReasonBaselineRewrite,
		RestatementReasonOperatorError,
	}
}
