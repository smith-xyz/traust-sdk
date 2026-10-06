// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// ReviewItemStatus — Queue state.
type ReviewItemStatus string

const (
	ReviewItemStatusPending   ReviewItemStatus = "pending"
	ReviewItemStatusConfirmed ReviewItemStatus = "confirmed"
	ReviewItemStatusRejected  ReviewItemStatus = "rejected"
)

// ReviewItemStatusValues returns all valid ReviewItemStatus values.
func ReviewItemStatusValues() []ReviewItemStatus {
	return []ReviewItemStatus{
		ReviewItemStatusPending,
		ReviewItemStatusConfirmed,
		ReviewItemStatusRejected,
	}
}
