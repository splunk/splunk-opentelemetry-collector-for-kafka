package common

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"testing"
	"text/template"
	"time"
)

const (
	HostEnvVar               = "CI_SPLUNK_HOST"
	UserEnvVar               = "CI_SPLUNK_USERNAME"
	HecToken                 = "CI_SPLUNK_HEC_TOKEN"
	PasswordEnvVar           = "CI_SPLUNK_PASSWORD"
	ManagementPortEnvVar     = "CI_SPLUNK_MGMT_PORT"
	KafkaBrokerAddressEnvVar = "CI_KAFKA_BROKER_ADDRESS"
	OTel_Binary              = "CI_OTEL_BINARY_FILE"
)

const (
	EventSearchQueryString = "| search "
	ConfigFilesDir         = "./testdata/configs"
	TestCaseDuration       = 60 * time.Second
	PerfTestCaseDuration   = 10 * time.Minute
	TestCaseTick           = 5 * time.Second
)

const (
	minIngestRate_num_1000000_bytes_600   = 13.829119
	minIngestRate_num_10000000_bytes_10   = 0.612903
	minIngestRate_num_10000000_bytes_100  = 5.324022
	minIngestRate_num_1000000_bytes_10000 = 37.327171
	minIngestRate_num_1000000_bytes_1000  = 18.374272
	minIngestRate_num_1000000_bytes_300   = 9.335118
	minIngestRateThreshold                = 0.9
)

// GetConfigVariable returns the value of the environment variable with the given name.
func GetConfigVariable(variableName string) string {
	envVariableName := ""
	switch variableName {
	case "HOST":
		envVariableName = HostEnvVar
	case "HEC_TOKEN":
		envVariableName = HecToken
	case "USER":
		envVariableName = UserEnvVar
	case "PASSWORD":
		envVariableName = PasswordEnvVar
	case "MANAGEMENT_PORT":
		envVariableName = ManagementPortEnvVar
	case "KAFKA_BROKER_ADDRESS":
		envVariableName = KafkaBrokerAddressEnvVar
	case "OTEL_BINARY_FILE":
		envVariableName = OTel_Binary
	}
	value := os.Getenv(envVariableName)

	if value != "" {
		return value
	}
	panic(envVariableName + " environment variable is not set")
}

func PrepareConfigFile(t *testing.T, configTemplateFile string, replacements map[string]any, configFilesDir string) string {
	cfgBytes, err := os.ReadFile(fmt.Sprintf("%s/%s", configFilesDir, configTemplateFile))
	require.NoError(t, err)
	tmpl, err := template.New("").Parse(string(cfgBytes))
	require.NoError(t, err)
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, replacements)
	require.NoError(t, err)
	var oTelKafkaConfig map[string]any
	err = yaml.Unmarshal(buf.Bytes(), &oTelKafkaConfig)
	require.NoError(t, err)

	//fmt.Printf("config: %v\n", oTelKafkaConfig) // debugging line to print the config
	// save the config to a file
	trimmedFileName := strings.TrimSuffix(configTemplateFile, ".tmpl")
	configFilePath := fmt.Sprintf("%s/%s", configFilesDir, trimmedFileName)
	err = os.WriteFile(configFilePath, buf.Bytes(), 0644)
	require.NoError(t, err, "Failed to write config file")
	// Print the path of the config file
	t.Logf("Config file created: %s\n", configFilePath)
	return trimmedFileName
}

func GetMinimumIngestRate() (float64, error) {
	numMsg := os.Getenv("NUM_MSG")
	recordSize := os.Getenv("RECORD_SIZE")

	switch {
	case numMsg == "1000000" && recordSize == "600":
		return minIngestRateThreshold * minIngestRate_num_1000000_bytes_600, nil
	case numMsg == "10000000" && recordSize == "10":
		return minIngestRateThreshold * minIngestRate_num_10000000_bytes_10, nil
	case numMsg == "10000000" && recordSize == "100":
		return minIngestRateThreshold * minIngestRate_num_10000000_bytes_100, nil
	case numMsg == "1000000" && recordSize == "10000":
		return minIngestRateThreshold * minIngestRate_num_1000000_bytes_10000, nil
	case numMsg == "1000000" && recordSize == "1000":
		return minIngestRateThreshold * minIngestRate_num_1000000_bytes_1000, nil
	case numMsg == "1000000" && recordSize == "300":
		return minIngestRateThreshold * minIngestRate_num_1000000_bytes_300, nil
	default:
		return 0, fmt.Errorf("unknown NUM_MSG (%s) and RECORD_SIZE (%s) combination", numMsg, recordSize)
	}
}
