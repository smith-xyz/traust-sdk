// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// IsolationCheckResult — Outcome of one isolation-review dimension or PQC readiness check. `na` means not applicable.
type IsolationCheckResult string

const (
	IsolationCheckResultYes     IsolationCheckResult = "yes"
	IsolationCheckResultPartial IsolationCheckResult = "partial"
	IsolationCheckResultNo      IsolationCheckResult = "no"
	IsolationCheckResultNa      IsolationCheckResult = "na"
)

// IsolationCheckResultValues returns all valid IsolationCheckResult values.
func IsolationCheckResultValues() []IsolationCheckResult {
	return []IsolationCheckResult{
		IsolationCheckResultYes,
		IsolationCheckResultPartial,
		IsolationCheckResultNo,
		IsolationCheckResultNa,
	}
}
