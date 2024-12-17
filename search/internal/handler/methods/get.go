package methods

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type get struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Location []string `json:"location"`
}

func HandleGet(data map[string]any) (map[string]any, error) {
	var g get
	if err := mapstructure.Decode(data, &g); err != nil {
		return nil, fmt.Errorf("handle get: failed to parse message: %w", err)
	}
	fmt.Println(g)
	return map[string]any{}, nil
}
