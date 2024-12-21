package methods

import (
	"fmt"

	repo "gitlab.crja72.ru/gospec/go21/go_final_project/internal/repository"
)

func HandleAdd(data map[string]any, rep repo.Repository) (repo.AddResp, error) {
	var a repo.AddReq
	a.Name = data["name"].(string)
	a.Description = data["description"].(string)
	a.Type = data["type"].(string)
	a.UserLogin = data["login"].(string)

	a.Location.City = data["location"].(map[string]any)["city"].(string)
	a.Location.Country = data["location"].(map[string]any)["country"].(string)
	a.Location.District = data["location"].(map[string]any)["district"].(string)

	if a.UserLogin == "" {
		return repo.AddResp{}, fmt.Errorf("handle add: failed to parse message: userlogin not provided")
	}
	fmt.Println(a)
	return rep.AddFind(a)
}
