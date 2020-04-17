package functions

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
)

// Read in the Credentials from command line or ENV variables
func ReadCredentials() (string, string, string) {
	//Read in Env Variable if they exist )
	username, exists := os.LookupEnv("USER")
	if !exists {
		fmt.Print("Enter Username: ")
		username, _ := terminal.ReadPassword(0)
		fmt.Printf("\n")
		strings.TrimSpace(string(username))
	}

	password, exists := os.LookupEnv("JIRAPW")
	if !exists {
		fmt.Printf("Password: ")
		password, _ := terminal.ReadPassword(0)
		fmt.Printf("\n")
		strings.TrimSpace(string(password))
	}

	slackToken, exists := os.LookupEnv("SLACKTOKENCOM")
	if !exists {
		fmt.Printf("Slack Token: ")
		slackToken, _ := terminal.ReadPassword(0)
		fmt.Printf("\n")
		strings.TrimSpace(string(slackToken))
	}

	return username, password, slackToken
}
