package types

type MetaData struct {
	SchemaVersion int    `json:"schema_version"` // версия структуры
	AppVersion    string `json:"app_version"`    // версия программы
}
