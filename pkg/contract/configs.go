package contract

// Represents request payload for creating/updating a config
type UpsertConfigRequest struct {
	Config
}

// Represents the response payload for a config
type GetConfigResponse struct {
	Config
}

// Represents a config
type Config struct {
	Name     string      `json:"name"`
	Metadata interface{} `json:"metadata"`
}
