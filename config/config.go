package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
)

const (
	jiraBaseUrlEnv   = "JIRA_BASE_URL"
	jiraAPITokenEnv  = "JIRA_API_TOKEN"
	jiraUserEmailEnv = "JIRA_USER_EMAIL"
)

var (
	readFile = ioutil.ReadFile
	getenv   = os.Getenv
)

// JiraConfig specifies the login information
type JiraConfig struct {
	BaseURL   string `json:"baseUrl"`
	APIToken  string `json:"token"`
	UserEmail string `json:"email"`
}

func ReadConfig() (JiraConfig, error) {
	config := JiraConfig{}
	homeDir := getenv("HOME")
	configPath := path.Join(homeDir, "jira", "config.yml")

	configYaml, err := readFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return config, err
		}

		// Fall back to env variables if the config file doesn't exist
		config.BaseURL = getenv(jiraBaseUrlEnv)
		config.APIToken = getenv(jiraAPITokenEnv)
		config.UserEmail = getenv(jiraUserEmailEnv)
	} else {
		yaml.Unmarshal(configYaml, &config)
	}

	emptyConfig := JiraConfig{}
	if config == emptyConfig {
		return config, errors.New("no config specified, please either use atlassian/gajira-login or provide the JIRA_ env variables specified in the readme")
	}

	return config, nil
}
