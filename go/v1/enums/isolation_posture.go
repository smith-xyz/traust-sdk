// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// IsolationPosture — Values of `overall` in isolation-review.schema.json.
type IsolationPosture string

const (
	IsolationPostureStrong      IsolationPosture = "strong"
	IsolationPostureAdequate    IsolationPosture = "adequate"
	IsolationPostureWeak        IsolationPosture = "weak"
	IsolationPostureCriticalGap IsolationPosture = "critical-gap"
)

// IsolationPostureValues returns all valid IsolationPosture values.
func IsolationPostureValues() []IsolationPosture {
	return []IsolationPosture{
		IsolationPostureStrong,
		IsolationPostureAdequate,
		IsolationPostureWeak,
		IsolationPostureCriticalGap,
	}
}
