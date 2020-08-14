package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/tomhjp/gh-action-jira-search/config"
)

var schemeRegex = regexp.MustCompile("https?")

type searchResult struct {
	Issues []struct {
		Key string `json:"key"`
	} `json:"issues"`
}

// FindIssueKeys queries Jira with some JQL and returns a slice of all the returned issue keys
func FindIssueKeys(config config.JiraConfig, jql string) ([]string, error) {
	url := generateURL(config.BaseURL, "/rest/api/3/search", url.Values{
		"jql":    {jql},
		"fields": {"summary"}, // Specify fields summary purely to minimise the size of all the unused fields in the response.
	})
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	if config.UserEmail != "" && config.APIToken != "" {
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", config.UserEmail, config.APIToken))))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API call failed (%d): %s", resp.StatusCode, string(bytes))
	}

	searchResult := searchResult{}
	json.Unmarshal(bytes, &searchResult)

	result := []string{}
	for _, issue := range searchResult.Issues {
		result = append(result, issue.Key)
	}

	return result, nil
}

func generateURL(baseURL string, path string, query url.Values) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	return fmt.Sprintf("%s%s?%s", baseURL, path, query.Encode())
}
