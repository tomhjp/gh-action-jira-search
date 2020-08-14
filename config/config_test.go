package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	configURL   = "https://config.local"
	configToken = "configSecret"
	configEmail = "user@config.local"
	envURL      = "https://env.local"
	envToken    = "envSecret"
	envEmail    = "user@env.local"
)

var (
	configContents = fmt.Sprintf(`baseUrl: %s
token: %s
email: %s`, configURL, configToken, configEmail)

	env = map[string]string{
		jiraBaseUrlEnv:   envURL,
		jiraAPITokenEnv:  envToken,
		jiraUserEmailEnv: envEmail,
	}
)

func TestReadConfig(t *testing.T) {
	envNotPopulated := func(string) string { return "" }
	envPopulated := func(s string) string { return env[s] }
	noConfigFile := func(s string) ([]byte, error) { return nil, os.ErrNotExist }
	configExists := func(s string) ([]byte, error) { return []byte(configContents), nil }

	for _, tc := range []struct {
		name               string
		readFileFunc       func(string) ([]byte, error)
		envFunc            func(string) string
		expectError        bool
		expectConfigValues bool
	}{
		{
			"No inputs, should error",
			noConfigFile,
			envNotPopulated,
			true,
			false,
		},
		{
			"Just config file",
			configExists,
			envNotPopulated,
			false,
			true,
		},
		{
			"Just env variables",
			noConfigFile,
			envPopulated,
			false,
			false,
		},
		{
			"Both, should use config file",
			configExists,
			envPopulated,
			false,
			true,
		},
		{
			"Permission denied error on reading file",
			func(string) ([]byte, error) { return nil, os.ErrPermission },
			envPopulated,
			true,
			false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			readFile = tc.readFileFunc
			getenv = tc.envFunc

			config, err := ReadConfig()
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tc.expectConfigValues {
					assert.Equal(t, configURL, config.BaseURL)
					assert.Equal(t, configToken, config.APIToken)
					assert.Equal(t, configEmail, config.UserEmail)
				} else {
					assert.Equal(t, envURL, config.BaseURL)
					assert.Equal(t, envToken, config.APIToken)
					assert.Equal(t, envEmail, config.UserEmail)
				}
			}
		})
	}
}
