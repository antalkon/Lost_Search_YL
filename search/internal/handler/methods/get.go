package methods

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	repo "gitlab.crja72.ru/gospec/go21/go_final_project/internal/repository"
)

func HandleGet(data map[string]any, rep repo.Repository) (repo.GetResp, error) {
	var g repo.GetReq
	if err := mapstructure.Decode(data, &g); err != nil {
		return repo.GetResp{}, fmt.Errorf("handle get: failed to parse message: %w", err)
	}
	fmt.Println(g)
	return rep.GetFind(g)
}
