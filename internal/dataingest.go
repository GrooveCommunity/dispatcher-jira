package internal

import (
	"log"
	"os"

	"encoding/json"

	"github.com/GrooveCommunity/glib-cloud-storage/gcp"
	"github.com/GrooveCommunity/glib-noc-event-structs/entity"
)

var bucketnames = os.Getenv("GCP_RULES_BUCKET")

func WriteRule(rule entity.Rule) {
	gcp.WriteObject(rule, bucketnames, rule.Name)
}

func GetRules() []entity.Rule {
	var rules []entity.Rule

	dataObjects := gcp.GetObjects(bucketnames)

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
