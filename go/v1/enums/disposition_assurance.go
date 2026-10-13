// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// DispositionAssurance — Level of confidence in a disposition. Indicates how a status was verified.
type DispositionAssurance string

const (
	DispositionAssuranceExecutionProven DispositionAssurance = "execution_proven"
	DispositionAssuranceHumanReviewed   DispositionAssurance = "human_reviewed"
	DispositionAssuranceMachineVerified DispositionAssurance = "machine_verified"
	DispositionAssuranceClaimed         DispositionAssurance = "claimed"
)

// DispositionAssuranceValues returns all valid DispositionAssurance values.
func DispositionAssuranceValues() []DispositionAssurance {
	return []DispositionAssurance{
		DispositionAssuranceExecutionProven,
		DispositionAssuranceHumanReviewed,
		DispositionAssuranceMachineVerified,
		DispositionAssuranceClaimed,
	}
}
