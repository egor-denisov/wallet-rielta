package config

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

var testConfigStr = `
app:
  name: "test-app"
  version: "1.0.0"
  countWorkers: 24
  timeout: 5s
  defaultBalance: 100

http:
  port: ":8080"
  timeout: 10s

postgres:
  poolMax: 2

rabbitmq:
  rpcServerExchange: "rpc_server"
  rpcClientExchange: "rpc_client"

logger:
  logLevel: "info"
`

var testEnvRequiredStr = `
PG_URL=test-url
RMQ_URL=test-url
`

var testEnvStr = `
APP_NAME=test-app
APP_VERSION=1.0.0
APP_WORKERS=24
APP_TIMEOUT=5s
APP_DEFAULT_BALANCE=100
HTTP_PORT=:8080
HTTP_TIMEOUT=5s
PG_POOL_MAX=2
PG_URL=test-url
RMQ_RPC_SERVER=rpc_server
RMQ_RPC_CLIENT=rpc_client
RMQ_URL=test-url
LOG_LEVEL=info`

func Test_MustLoadPath_ExistentPath(t *testing.T) {
	for _, test := range testsMustLoadPath {
		t.Run(test.name, func(t *testing.T) {
			// temp config.yml file
			tempFileConfig, err := os.CreateTemp("", "config-*.yml")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tempFileConfig.Name())

			_, err = tempFileConfig.WriteString(test.configFile)
			if err != nil {
				t.Fatal(err)
			}
			// temp .env file
			tempFileEnv, err := os.CreateTemp("", "*.env")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tempFileEnv.Name())

			_, err = tempFileEnv.WriteString(test.envFile)
			if err != nil {
				t.Fatal(err)
			}

			config := MustLoadPath(tempFileConfig.Name(), tempFileEnv.Name())
			assert.Equal(t, test.expectedConfig, config)
		})
	}
}

var testsMustLoadPath = []struct {
	name           string
	configFile     string
	envFile        string
	expectedConfig *Config
}{
	{
		name:       "Only .env",
		configFile: "something: else",
		envFile:    testEnvStr,
		expectedConfig: &Config{
			App: App{
				Name:           "test-app",
				Version:        "1.0.0",
				CountWorkers:   24,
				Timeout:        5 * time.Second,
				DefaultBalance: 100,
			},
			HTTP: HTTP{
				Port:    ":8080",
				Timeout: 5 * time.Second,
			},
			PG: PG{
				PoolMax: 2,
				URL:     "test-url",
			},
			RMQ: RMQ{
				ServerExchange: "rpc_server",
				ClientExchange: "rpc_client",
				URL:            "test-url",
			},
			Log: Log{
				Level: "info",
			},
		},
	},
	{
		name:       "In .env only required",
		configFile: testConfigStr,
		envFile:    testEnvRequiredStr,
		expectedConfig: &Config{
			App: App{
				Name:           "test-app",
				Version:        "1.0.0",
				CountWorkers:   24,
				Timeout:        5 * time.Second,
				DefaultBalance: 100,
			},
			HTTP: HTTP{
				Port:    ":8080",
				Timeout: 5 * time.Second,
			},
			PG: PG{
				PoolMax: 2,
				URL:     "test-url",
			},
			RMQ: RMQ{
				ServerExchange: "rpc_server",
				ClientExchange: "rpc_client",
				URL:            "test-url",
			},
			Log: Log{
				Level: "info",
			},
		},
	},
}

func Test_MustLoadPath_NonExistentPath(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	MustLoadPath("non_existent_config.yml", "non_existent_env.env")
}

var testsFetchConfigPath = []struct {
	name      string
	argsValue []string
	envValue  string
	expected  string
}{
	{
		name:      "Not field",
		argsValue: []string{"cmd", ""},
		envValue:  "",
		expected:  _defaultConfigPath,
	},
	{
		name:      "Ok - from environment",
		argsValue: []string{"cmd", ""},
		envValue:  "./test_config2.yml",
		expected:  "./test_config2.yml",
	},
	{
		name:      "Ok - from flag",
		argsValue: []string{"cmd", "-config", "./test_config.yml"},
		envValue:  "",
		expected:  "./test_config.yml",
	},
}

func Test_fetchConfigPath(t *testing.T) {
	for _, test := range testsFetchConfigPath {
		t.Run(test.name, func(t *testing.T) {
			os.Args = test.argsValue
			t.Setenv("CONFIG_PATH", test.envValue)

			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			configPath := fetchConfigPath()

			assert.Equal(t, test.expected, configPath)
		})
	}
}
