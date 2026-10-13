// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// BenchmarkTargetEmbargo — Values of `embargo` in benchmark-target.schema.json.
type BenchmarkTargetEmbargo string

const (
	BenchmarkTargetEmbargoInternal BenchmarkTargetEmbargo = "internal"
	BenchmarkTargetEmbargoPublic   BenchmarkTargetEmbargo = "public"
)

// BenchmarkTargetEmbargoValues returns all valid BenchmarkTargetEmbargo values.
func BenchmarkTargetEmbargoValues() []BenchmarkTargetEmbargo {
	return []BenchmarkTargetEmbargo{
		BenchmarkTargetEmbargoInternal,
		BenchmarkTargetEmbargoPublic,
	}
}
