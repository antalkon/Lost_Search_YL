package methods

import (
	"fmt"

	repo "gitlab.crja72.ru/gospec/go21/go_final_project/internal/repository"
)

func HandleRespond(data map[string]any, rep repo.Repository) (repo.RespondResp, error) {
	var r repo.RespondReq
	var ok bool
	fmt.Println(data)

	r.FindUUID, ok = data["find_uuid"].(string)
	if !ok {
		return repo.RespondResp{}, fmt.Errorf("handle respond: failed to parse message: no uuid provided")
	}
	fmt.Println(r)
	return rep.RespondToFind(r)
}
