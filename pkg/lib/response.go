package lib

import (
	"encoding/json"
	"io"

	"go.uber.org/zap"
)

func WriteResponseJSON(w io.Writer, response interface{}) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	err := enc.Encode(response)
	if err != nil {
		zap.S().Errorf("write JSON response: %v", err)
	}
}
