// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// ForkVisibility — Values of `visibility` in remediation.schema.json.
type ForkVisibility string

const (
	ForkVisibilityPrivate  ForkVisibility = "private"
	ForkVisibilityInternal ForkVisibility = "internal"
	ForkVisibilityPublic   ForkVisibility = "public"
)

// ForkVisibilityValues returns all valid ForkVisibility values.
func ForkVisibilityValues() []ForkVisibility {
	return []ForkVisibility{
		ForkVisibilityPrivate,
		ForkVisibilityInternal,
		ForkVisibilityPublic,
	}
}
