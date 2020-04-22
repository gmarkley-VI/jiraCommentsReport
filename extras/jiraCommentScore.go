package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/montanaflynn/stats"
	"github.com/openshift/gmarkley-VI/jiraSosRepot/functions"
	"gopkg.in/jdkato/prose.v2"
	"log"
	"strings"
)

func scoreComment(comment string) float64 {
	doc, err := prose.NewDocument(comment)
	if err != nil {
		fmt.Printf("%s\n", err)
		return 0
	}

	// Iterate over the doc's tokens and count:
	var nounCount float64
	var verbCount float64
	for _, tok := range doc.Tokens() {
		if strings.Contains(tok.Tag, "NN") {
			nounCount++
		}
		if strings.Contains(tok.Tag, "VB") {
			verbCount++
		}
	}

	// count the doc's named-entities:
	var entitiesCount float64 = float64(len(doc.Entities()))

	// count the doc's sentences:
	var sentenceCount float64 = float64(len(doc.Sentences()))

	// calculate score
	var score = entitiesCount*10 + sentenceCount*2.5 + nounCount + verbCount
	if score >= 100 {
		score = 99
	}
	fmt.Printf("\n------------------------------------------------------\n%v\n", comment)
	//fmt.Printf("Entities: %v, Sentances: %v, Nouns: %v, Verb: %v, Score: %v\n", entitiesCount, sentenceCount, nounCount, verbCount, score)
	fmt.Printf("Score: %v\n", score)
	return score
}

func main() {
	//Setup
	jiraURL := "https://issues.redhat.com"
	username, password, _ := functions.ReadCredentials()
	var jiraJQL [1][2]string
	//jiraJQL[0][0] = "project = WINC AND status in (\"In Progress\", \"Code Review\")AND(sprint in openSprints())"
	jiraJQL[0][0] = "project = WINC AND issuetype = Story AND fixVersion = \"OpenShift 4.5\""

	//Create the client
	client, _ := functions.CreatTheClient(username, password, jiraURL)

	var data []float64

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

			if len(u.RenderedFields.Comments.Comments) >= 1 {
				c := u.RenderedFields.Comments.Comments[len(u.RenderedFields.Comments.Comments)-1]
				//fmt.Printf("\nScore: %v\nComment: %s\n------------------------------------------\n", scoreComment(c.Body), c.Body)
				data = append(data, scoreComment(c.Body))
			}
		}
	}
	var min, _ = stats.Min(data)
	var max, _ = stats.Max(data)
	var mean, _ = stats.Median(data)
	var mode, _ = stats.Mode(data)
	fmt.Printf("\nMin: %v, Avarage: %v, Mode: %v, Max: %v", min, mean, mode, max)
}
