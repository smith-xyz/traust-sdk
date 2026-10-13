// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// ThreatLikelihood — How often this threat is expected to be driven, on an ordinal scale.
type ThreatLikelihood string

const (
	ThreatLikelihoodVeryRare      ThreatLikelihood = "very_rare"
	ThreatLikelihoodRare          ThreatLikelihood = "rare"
	ThreatLikelihoodPossible      ThreatLikelihood = "possible"
	ThreatLikelihoodLikely        ThreatLikelihood = "likely"
	ThreatLikelihoodAlmostCertain ThreatLikelihood = "almost_certain"
)

// ThreatLikelihoodValues returns all valid ThreatLikelihood values.
func ThreatLikelihoodValues() []ThreatLikelihood {
	return []ThreatLikelihood{
		ThreatLikelihoodVeryRare,
		ThreatLikelihoodRare,
		ThreatLikelihoodPossible,
		ThreatLikelihoodLikely,
		ThreatLikelihoodAlmostCertain,
	}
}
