package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/openshift/gmarkley-VI/jiraSosRepot/functions"
	"log"
	"strings"
)

func main() {
	jiraURL := "https://issues.redhat.com"
	username, password := functions.ReadCredentials()

	var jiraJQL [1][2]string
	jiraJQL[0][0] = "project = WINC AND status in (\"In Progress\", \"Code Review\")AND(sprint in openSprints())"
	jiraJQL[0][1] = "--Current Winc Work Items--"

	//Create the client
	client, _ := functions.CreatTheClient(username, password, jiraURL)

	//Loop over the jiraJQL array and Request the issue objects
	for z := 0; z < len(jiraJQL); z++ {

		var issues []jira.Issue

		// append the jira issues to []jira.Issue
		appendFunc := func(i jira.Issue) (err error) {
			issues = append(issues, i)
			return err
		}

		// SearchPages will page through results and pass each issue to appendFunc taken from the Jira Example implementation
		// In this example, we'll search for all the issues with the provided JQL filter and Print the header that goes with it.
		err := client.Issue.SearchPages(fmt.Sprintf(`%s`, jiraJQL[z][0]), nil, appendFunc)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\n%s\n", jiraJQL[z][1])

		for _, i := range issues {
			options := &jira.GetQueryOptions{Expand: "renderedFields"}
			u, _, err := client.Issue.Get(i.Key, options)

			if err != nil {
				fmt.Printf("\n==> error: %v\n", err)
				return
			}
			c := u.RenderedFields.Comments.Comments[len(u.RenderedFields.Comments.Comments)-1]
			if strings.Contains(c.Updated, "days ago") || strings.Contains(c.Updated, "Yesterday") {
				fmt.Printf("%s Please comment/update - Last was %+v - ", i.Fields.Assignee.DisplayName, c.Updated)
				fmt.Printf("%s/browse/%s \n", strings.TrimSpace(jiraURL), i.Key)
			}
		}
	}
}
