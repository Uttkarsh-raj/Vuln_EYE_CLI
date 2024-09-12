package models

type ManifestData struct {
	ExportedActivityFound bool
	ExportedActivityTrue  bool
	Permissions           []string
	DebuggableFlag        string
	AllowBackupFlag       string
	IntentFilters         []string
}
