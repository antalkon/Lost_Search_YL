package methods

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type respond struct {
	FindUUID string `json:"find_uuid"`
}

func HandleRespond(data map[string]any) (map[string]any, error) {
	var r respond
	if err := mapstructure.Decode(data, &r); err != nil {
		return nil, fmt.Errorf("handle respond: failed to parse message: %w", err)
	}
	fmt.Println(r)
	return map[string]any{}, nil
}
