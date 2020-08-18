package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/tomhjp/gh-action-jira-search/config"
	"github.com/tomhjp/gh-action-jira-search/jira"
)

func main() {
	err := search()
	if err != nil {
		log.Fatal(err)
	}
}

func search() error {
	jql := os.Getenv("INPUT_JQL")
	if jql == "" {
		return errors.New("no jql query provided as input")
	}
	config, err := config.ReadConfig()
	if err != nil {
		return err
	}

	issueKeys, err := jira.FindIssueKeys(config, jql)
	if err != nil {
		return err
	}
	if len(issueKeys) == 0 {
		fmt.Println("Successfully queried API but did not find any issues")
		return nil
	} else if len(issueKeys) > 1 {
		return errors.New("jql does not uniquely identify an issue")
	}

	key := issueKeys[0]
	fmt.Printf("Found issue %s\n", key)

	// Special format log line to set output for the action.
	// See https://docs.github.com/en/actions/creating-actions/metadata-syntax-for-github-actions#outputs-for-composite-run-steps-actions.
	fmt.Printf("::set-output name=key::%s\n", key)

	return nil
}
