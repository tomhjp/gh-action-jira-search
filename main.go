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
		return errors.New("no issues found for jql query")
	} else if len(issueKeys) > 1 {
		return errors.New("jql does not uniquely identify an issue")
	}

	fmt.Printf(issueKeys[0])

	return nil
}
