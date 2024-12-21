package methods

import (
	"fmt"

	repo "gitlab.crja72.ru/gospec/go21/go_final_project/internal/repository"
)

func HandleGet(data map[string]any, rep repo.Repository) (repo.GetResp, error) {
	var g repo.GetReq
	g.Name = data["name"].(string)
	g.Type = data["type"].(string)

	g.Location.City = data["location"].(map[string]any)["city"].(string)
	g.Location.Country = data["location"].(map[string]any)["country"].(string)
	g.Location.District = data["location"].(map[string]any)["district"].(string)

	fmt.Println(g)
	return rep.GetFind(g)
}
