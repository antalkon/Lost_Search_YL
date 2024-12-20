package methods

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	repo "gitlab.crja72.ru/gospec/go21/go_final_project/internal/repository"
)

func HandleAdd(data map[string]any, rep repo.Repository) (repo.AddResp, error) {
	var a repo.AddReq
	if err := mapstructure.Decode(data, &a); err != nil {
		return repo.AddResp{}, fmt.Errorf("handle add: failed to parse message: %w", err)
	}
	fmt.Println(a)
	return rep.AddFind(a)
}
