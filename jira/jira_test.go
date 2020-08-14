package jira

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/tomhjp/gh-action-jira-search/config"

	"github.com/stretchr/testify/require"
)

func TestGenerateURL(t *testing.T) {
	url := generateURL("https://org.atlassian.net/", "/rest/api/3/search", url.Values{
		"jql": {"project = FOO"},
	})
	require.Equal(t, "https://org.atlassian.net/rest/api/3/search?jql=project+%3D+FOO", url)
}

func TestFindIssueKeys(t *testing.T) {
	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(`{
			"expand": "names,schema",
			"startAt": 0,
			"maxResults": 50,
			"total": 1,
			"issues": [
			  {
				"expand": "operations,versionedRepresentations,editmeta,changelog,renderedFields",
				"id": "27438",
				"key": "FOO-23",
				"fields": {
				  "summary": "Fix the foo system"
				}
			  },
			  {
				"expand": "operations,versionedRepresentations,editmeta,changelog,renderedFields",
				"id": "27438",
				"key": "FOO-24",
				"fields": {
				  "summary": "Fix the bar system"
				}
			  }
			]
		  }`))
	}))
	defer testServer.Close()

	config := config.JiraConfig{
		BaseURL:   testServer.URL,
		APIToken:  "supersecretvalue",
		UserEmail: "user@example.com",
	}
	keys, err := FindIssueKeys(config, `project = "VAULT""`)
	require.NoError(t, err)
	require.Len(t, keys, 2)
	require.Equal(t, []string{"FOO-23", "FOO-24"}, keys)
}