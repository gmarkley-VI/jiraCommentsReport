package functions

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
)

// Read in the Credentials from command line.
func ReadCredentials() (string, string) {
	fmt.Print("Enter Username: ")
	username, _ := terminal.ReadPassword(0)

	fmt.Printf("\nPassword: ")
	password, _ := terminal.ReadPassword(0)
	fmt.Printf("\n")

	return strings.TrimSpace(string(username)), strings.TrimSpace(string(password))
}
