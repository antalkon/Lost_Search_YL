package methods

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	repo "gitlab.crja72.ru/gospec/go21/go_final_project/internal/repository"
)

func HandleRespond(data map[string]any, rep repo.Repository) (repo.RespondResp, error) {
	var r repo.RespondReq
	if err := mapstructure.Decode(data, &r); err != nil {
		return repo.RespondResp{}, fmt.Errorf("handle respond: failed to parse message: %w", err)
	}
	fmt.Println(r)
	return rep.RespondToFind(r)
}
