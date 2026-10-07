// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// PqcPathClass — Values of `path_class` in pqc-facts.schema.json.
type PqcPathClass string

const (
	PqcPathClassFirstParty PqcPathClass = "first_party"
	PqcPathClassVendor     PqcPathClass = "vendor"
	PqcPathClassTestDocs   PqcPathClass = "test_docs"
)

// PqcPathClassValues returns all valid PqcPathClass values.
func PqcPathClassValues() []PqcPathClass {
	return []PqcPathClass{
		PqcPathClassFirstParty,
		PqcPathClassVendor,
		PqcPathClassTestDocs,
	}
}
