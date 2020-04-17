package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/openshift/gmarkley-VI/jiraSosRepot/functions"
	"github.com/slack-go/slack"
	"log"
	"strings"
)

func exportConsole(key string, output string) {
	fmt.Printf("ID - %s - %s\n", key, output)
}

func exportJira(client *jira.Client, id string, key string, output string) *jira.Comment {
	com := jira.Comment{
		ID:           id,
		Self:         "",
		Name:         "",
		Author:       jira.User{},
		Body:         output,
		UpdateAuthor: jira.User{},
		Updated:      "",
		Created:      "",
		Visibility:   jira.CommentVisibility{},
	}
	commentOUT, _, err := client.Issue.AddComment(key, &com)
	if err != nil {
		panic(err)
	}
	return commentOUT
}

func exportSlack(token string, key string, output string) {
	api := slack.New(token)

	// Production channelID, timestamp, err := api.PostMessage("#forum-winc", slack.MsgOptionText(output, false))
	channelID, timestamp, err := api.PostMessage("#forum-winc", slack.MsgOptionText(output, false))
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("\"%s\" - Message successfully sent to channel %s at %s\n", output, channelID, timestamp)
}

func main() {
	//Setup
	jiraURL := "https://issues.redhat.com"
	username, password, slackToken := functions.ReadCredentials()
	var jiraJQL [1][2]string
	//jiraJQL[0][0] = "project = WINC AND status in (\"In Progress\", \"Code Review\")AND(sprint in openSprints())"
	jiraJQL[0][0] = "project = WINC"

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

		for _, i := range issues {
			options := &jira.GetQueryOptions{Expand: "renderedFields"}
			u, _, err := client.Issue.Get(i.Key, options)
			if err != nil {
				fmt.Printf("\n==> error: %v\n", err)
				return
			}
			var name = i.Fields.Assignee.DisplayName

			if len(u.RenderedFields.Comments.Comments) >= 1 {

				c := u.RenderedFields.Comments.Comments[len(u.RenderedFields.Comments.Comments)-1]
				if strings.Contains(c.Updated, "days ago") {
					commentString := fmt.Sprintf("%v Please comment/update - Last update was %+v", name, c.Updated)
					//exportJira(client, i.ID, i.Key, commentString)
					exportConsole(i.Key, commentString)
					exportSlack(slackToken, i.Key, commentString)
				}
			} else {
				commentString := fmt.Sprintf("%v Please add a comment.", name)
				//exportJira(client, i.ID, i.Key, commentString)
				exportConsole(i.Key, commentString)
				exportSlack(slackToken, i.Key, commentString)
			}
		}
	}
}
