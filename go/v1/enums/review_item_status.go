// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

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
