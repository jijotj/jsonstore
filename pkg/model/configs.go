package model

// Represents a config object in persistence
type Config struct {
	Name     string      `json:"name"`
	Metadata interface{} `json:"metadata"`
}
