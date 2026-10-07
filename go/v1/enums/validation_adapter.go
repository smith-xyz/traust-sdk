// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// ValidationAdapter — Values of `adapter` in validation.schema.json.
type ValidationAdapter string

const (
	ValidationAdapterK8s       ValidationAdapter = "k8s"
	ValidationAdapterContainer ValidationAdapter = "container"
	ValidationAdapterWasm      ValidationAdapter = "wasm"
)

// ValidationAdapterValues returns all valid ValidationAdapter values.
func ValidationAdapterValues() []ValidationAdapter {
	return []ValidationAdapter{
		ValidationAdapterK8s,
		ValidationAdapterContainer,
		ValidationAdapterWasm,
	}
}
