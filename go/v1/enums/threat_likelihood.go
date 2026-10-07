// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

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
