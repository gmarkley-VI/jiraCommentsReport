package functions

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"strings"
)

func CreatTheClient(username string, password string, jiraURL string) (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	//Create the client
	client, err := jira.NewClient(tp.Client(), strings.TrimSpace(jiraURL))
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
	}
	return client, nil
}
