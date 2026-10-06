// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// FileChangeType — Values of `change_type` in remediation.schema.json.
type FileChangeType string

const (
	FileChangeTypeModified FileChangeType = "modified"
	FileChangeTypeAdded    FileChangeType = "added"
	FileChangeTypeDeleted  FileChangeType = "deleted"
	FileChangeTypeRenamed  FileChangeType = "renamed"
)

// FileChangeTypeValues returns all valid FileChangeType values.
func FileChangeTypeValues() []FileChangeType {
	return []FileChangeType{
		FileChangeTypeModified,
		FileChangeTypeAdded,
		FileChangeTypeDeleted,
		FileChangeTypeRenamed,
	}
}
