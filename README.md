# gh-action-jira-search

Run JQL in a [GitHub action](https://docs.github.com/en/actions) to find a specific Jira issue key.

## Authentication

To provide a URL and credentials you can use the [`gajira-login`](https://github.com/atlassian/gajira-login) action, which will write a config file this action can read.
Alternatively, you can set some environment variables:

- `JIRA_BASE_URL` - e.g. `https://my-org.atlassian.net`. The URL for your Jira instance.
- `JIRA_API_TOKEN` - e.g. `iaJGSyaXqn95kqYvq3rcEGu884TCbMkU`. An access token.
- `JIRA_USER_EMAIL` - e.g. `user@example.com`. The email address for the access token.

## Inputs

- `jql` - JQL query that returns at most 1 issue

## Outputs

The action will exit with a zero exit code unless it encounters any errors or finds more than 1 issue.

- `issue` - The issue key found, e.g. TEST-23. Empty if none.

## Examples

Using `atlassian/gajira-login` and [GitHub secrets](https://docs.github.com/en/actions/configuring-and-managing-workflows/creating-and-storing-encrypted-secrets) for authentication:

```yaml
- name: Login
  uses: atlassian/gajira-login@v2.0.0
  env:
    JIRA_BASE_URL: ${{ secrets.JIRA_BASE_URL }}
    JIRA_USER_EMAIL: ${{ secrets.JIRA_USER_EMAIL }}
    JIRA_API_TOKEN: ${{ secrets.JIRA_API_TOKEN }}

- name: Search
  id: search
  uses: tomhjp/gh-action-jira-search@v0.1.0
  with:
    jql: 'key = TEST-23'

- name: Log
  run: echo "Found issue ${{ steps.search.outputs.issue }}"
```

Using environment variables for authentication:

```yaml
- name: Search
  id: search
  uses: tomhjp/gh-action-jira-search@v0.1.0
  with:
    jql: 'key = TEST-23'
  env:
    JIRA_BASE_URL: ${{ secrets.JIRA_BASE_URL }}
    JIRA_USER_EMAIL: ${{ secrets.JIRA_USER_EMAIL }}
    JIRA_API_TOKEN: ${{ secrets.JIRA_API_TOKEN }}
```
