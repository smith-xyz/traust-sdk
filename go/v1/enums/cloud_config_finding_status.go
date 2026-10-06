// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// CloudConfigFindingStatus — Values of `status` in cloud-config-audit.schema.json.
type CloudConfigFindingStatus string

const (
	CloudConfigFindingStatusConfirmed   CloudConfigFindingStatus = "confirmed"
	CloudConfigFindingStatusSuppressed  CloudConfigFindingStatus = "suppressed"
	CloudConfigFindingStatusNeedsReview CloudConfigFindingStatus = "needs_review"
)

// CloudConfigFindingStatusValues returns all valid CloudConfigFindingStatus values.
func CloudConfigFindingStatusValues() []CloudConfigFindingStatus {
	return []CloudConfigFindingStatus{
		CloudConfigFindingStatusConfirmed,
		CloudConfigFindingStatusSuppressed,
		CloudConfigFindingStatusNeedsReview,
	}
}
