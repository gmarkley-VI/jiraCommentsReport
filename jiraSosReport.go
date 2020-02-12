package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/openshift/gmarkley-VI/jiraSosRepot/functions"
	"log"
)

func main() {
	jiraURL := "https://issues.redhat.com"
	username, password := functions.ReadCredentials()

	var jiraJQL [3][2]string
	jiraJQL[0][0] = "project = WINC AND (resolved >= -7d OR (status in (Done, Pending) AND sprint in openSprints())) AND priority in (Blocker, Critical, High) ORDER BY priority DESC"
	jiraJQL[0][1] = "--Completed\\Completing Last Week--"
	jiraJQL[1][0] = "project = WINC AND (status in (\"In Progress\", \"Code Review\") AND sprint in openSprints()) AND priority in (Blocker, Critical, High) ORDER BY priority DESC"
	jiraJQL[1][1] = "--Currently Active--"
	jiraJQL[2][0] = "project = WINC AND (status in (\"To Do\") AND sprint in openSprints()) AND priority in (Blocker, Critical, High) ORDER BY priority DESC"
	jiraJQL[2][1] = "--Remaining in Sprint--"

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
			//fmt.Printf("%s: %s - https://issues.redhat.com/browse/%s\n", i.Fields.Type.Name, i.Fields.Summary, i.Key)
			fmt.Printf("%s\n", i.Fields.Summary)
		}
	}

}
