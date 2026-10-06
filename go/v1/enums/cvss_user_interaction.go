// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// CvssUserInteraction — Values of `user_interaction` in validation.schema.json.
type CvssUserInteraction string

const (
	CvssUserInteractionNone     CvssUserInteraction = "none"
	CvssUserInteractionRequired CvssUserInteraction = "required"
)

// CvssUserInteractionValues returns all valid CvssUserInteraction values.
func CvssUserInteractionValues() []CvssUserInteraction {
	return []CvssUserInteraction{
		CvssUserInteractionNone,
		CvssUserInteractionRequired,
	}
}
