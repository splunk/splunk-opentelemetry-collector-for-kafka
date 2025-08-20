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
	ConfigFilesDir         = "../testdata/configs"
	TestCaseDuration       = 30 * time.Second
	TestCaseTick           = 5 * time.Second
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
	fmt.Printf("Config file created: %s\n", configFilePath)
	return trimmedFileName
}
