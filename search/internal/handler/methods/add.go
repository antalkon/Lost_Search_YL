package methods

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"gitlab.crja72.ru/gospec/go21/go_final_project/internal/repository"
)

type add struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Location    []string `json:"location"`
}

func HandleAdd(data map[string]any, repo repository.Repository) (map[string]any, error) {
	var a add
	if err := mapstructure.Decode(data, &a); err != nil {
		return nil, fmt.Errorf("handle add: failed to parse message: %w", err)
	}
	repo.CreateFind()
	fmt.Println(a)
	return map[string]any{}, nil
}
