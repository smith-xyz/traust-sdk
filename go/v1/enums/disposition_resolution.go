// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// DispositionResolution — Finding remediation resolution state. Tracks progress from open through risk acceptance or regression.
type DispositionResolution string

const (
	DispositionResolutionOpen                 DispositionResolution = "open"
	DispositionResolutionFixInProgress        DispositionResolution = "fix_in_progress"
	DispositionResolutionResolved             DispositionResolution = "resolved"
	DispositionResolutionPartiallyResolved    DispositionResolution = "partially_resolved"
	DispositionResolutionRiskAccepted         DispositionResolution = "risk_accepted"
	DispositionResolutionRegressionIntroduced DispositionResolution = "regression_introduced"
)

// DispositionResolutionValues returns all valid DispositionResolution values.
func DispositionResolutionValues() []DispositionResolution {
	return []DispositionResolution{
		DispositionResolutionOpen,
		DispositionResolutionFixInProgress,
		DispositionResolutionResolved,
		DispositionResolutionPartiallyResolved,
		DispositionResolutionRiskAccepted,
		DispositionResolutionRegressionIntroduced,
	}
}
