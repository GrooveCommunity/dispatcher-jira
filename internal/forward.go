package internal

import (
	"strings"

	"bytes"
	//	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/GrooveCommunity/glib-cloud-storage/gcp"
	"github.com/GrooveCommunity/glib-noc-event-structs/entity"
	"github.com/andygrunwald/go-jira"
)

func ForwardIssue(jiraIssue entity.JiraIssue, username, token, endpoint string) {

	rules := GetRules()

	tp := jira.BasicAuthTransport{
		Username: username, //usuário do jira
		Password: token,    //token de api
	}

	client, err := jira.NewClient(tp.Client(), strings.TrimSpace(endpoint))
	if err != nil {
		panic("\nError:" + err.Error())
	}

	for _, rule := range rules {
		if !validateRule(jiraIssue.CustomFields, rule.Forward.Input.Fields) {
			continue
		}

		//regra considera se é para ignorar o anexo, caso não seja a regra irá validar os cenários onde o anexo precise existir e o cenário onde o anexo não deva ser informado
		if !rule.Forward.Input.IgnoreAttachment && ((rule.Forward.Input.HasAttachment && len(jiraIssue.Attachment) == 0) || (!rule.Forward.Input.HasAttachment && len(jiraIssue.Attachment) > 0)) {
			continue
		}

		//verifica se foi definido alguma regra para validação de conteúdo na description
		if len(rule.Forward.Input.Contents) > 0 {
			//Verificar os conteúdos informados na regra
			for _, content := range rule.Forward.Input.Contents {
				//Valida se existe o conteúdo informado na regra no campo description
				if !strings.Contains(jiraIssue.Description, content) {
					//A regra não é valida se não existir o conteúdo informado no campo description
					continue
				}
			}
		}

		updateStatusIssue(client, jiraIssue.IssueID, "Analisar - SD")
		updateStatusIssue(client, jiraIssue.IssueID, "Acionar Squad")

		updateIssueCustomField(entity.JiraForwarded{Issue: jiraIssue, Rule: rule})
	}
}

func updateStatusIssue(client *jira.Client, issueID, status string) {
	fmt.Println("Issue ID:" + issueID)

	var transitionID string
	possibleTransitions, _, err := client.Issue.GetTransitions(issueID)

	if err != nil {
		panic("\nError:" + err.Error())
	}

	for _, transitions := range possibleTransitions {
		if transitions.Name == status {
			transitionID = transitions.ID
			break
		}
	}

	_, errorTransition := client.Issue.DoTransition(issueID, transitionID)

	if errorTransition != nil {
		panic("\nError:" + errorTransition.Error())
	}

	fmt.Println("Status atualizado para " + status)

}

func updateIssueCustomField(jiraForwared entity.JiraForwarded) {
	host := os.Getenv("JIRA_ENDPOINT") + "/rest/api/2/issue/" + jiraForwared.Issue.Key

	data := "{\"fields\": {\"" + jiraForwared.Rule.Forward.Output.CustomFieldID + "\":{\"value\":\"" + jiraForwared.Rule.Forward.Output.CustomFieldValue + "\"}}}"

	req, err := http.NewRequest(http.MethodPut, host, bytes.NewBuffer([]byte(data)))
	req.SetBasicAuth(os.Getenv("JIRA_USERNAME"), os.Getenv("JIRA_TOKENAPI"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	gcp.WriteObject(jiraForwared, "forwarded-calls", jiraForwared.Issue.Key)

	fmt.Println("Issue " + jiraForwared.Issue.Key + " atualizada!")
}

func validateRule(customFields []entity.CustomField, fields []entity.Field) bool {
	ruleState := true

	for _, field := range fields {
		ruleFieldState := false
		for _, customField := range customFields {
			if field.ID == customField.CustomID && field.Value == customField.Value {
				ruleFieldState = true
			}
		}

		//Caso os requisitos para a regra não sejam atendidas, retorna falso
		ruleState = ruleState && ruleFieldState
	}

	return ruleState
}
