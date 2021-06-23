package internal

import (
	"log"

	"encoding/json"

	gcp "github.com/GrooveCommunity/glib-cloud-storage/gcp"
	"github.com/GrooveCommunity/glib-noc-event-structs/entity"
)

func WriteRule(rule entity.Rule) {
	gcp.WriteObject(rule, "rules-dispatcher", rule.Name)

	//Atualiza a estrutura de rules
	UpdateRules(rule)
}

func GetRules() []entity.Rule {
	var rules []entity.Rule

	dataObjects := gcp.GetObjects("rules-dispatcher")

	for _, b := range dataObjects {
		var rule entity.Rule
		errUnmarsh := json.Unmarshal(b, &rule)

		if errUnmarsh != nil {
			log.Fatal("Erro no unmarshal\n", errUnmarsh.Error())
		}

		rules = append(rules, rule)
	}

	return rules
}
