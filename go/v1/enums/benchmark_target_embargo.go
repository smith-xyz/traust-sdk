// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

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
