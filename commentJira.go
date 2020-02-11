package main

import (
	"bufio"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/openshift/gmarkley-VI/jiraSosRepot/functions"
	"os"
	"strings"
)

func readComment() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter JIRA ID Number: ")
	idNum, _ := reader.ReadString('\n')

	fmt.Print("Enter Comment: ")
	comment, _ := reader.ReadString('\n')

	return strings.TrimSpace(idNum), comment
}

func main() {
	jiraURL := "https://issues.redhat.com"
	username, password := functions.ReadCredentials()

	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	//Create the client
	client, err := jira.NewClient(tp.Client(), strings.TrimSpace(jiraURL))
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	id, commentString := readComment()

	i := jira.Comment{
		ID:           id,
		Self:         "",
		Name:         "",
		Author:       jira.User{},
		Body:         commentString,
		UpdateAuthor: jira.User{},
		Updated:      "",
		Created:      "",
		Visibility:   jira.CommentVisibility{},
	}

	commentOUT, _, err := client.Issue.AddComment(id, &i)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ID - %s \n Body - %+v\n", commentOUT.ID, commentOUT.Body)
}
