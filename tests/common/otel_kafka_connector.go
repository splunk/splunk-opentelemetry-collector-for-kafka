package common

import (
	"fmt"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func StopOTelKafkaConnector(t *testing.T, cmd *exec.Cmd) {
	// Kill the process
	err := cmd.Process.Signal(syscall.SIGTERM)

	require.NoError(t, err)

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(10 * time.Second):
		// Timeout: process did not exit, kill it
		_ = cmd.Process.Kill()
		t.Logf("Process force killed after timeout.")
	case err := <-done:
		require.NoError(t, err)
		t.Logf("Process exited gracefully.")
	}
}

func StartOTelKafkaConnector(t *testing.T, configName string, configFilesDir string, featureGate ...string) *exec.Cmd {
	// Start the process
	var featureGatesArgs []string
	if len(featureGate) > 0 {
		featureGatesArgs = []string{"--feature-gates=" + featureGate[0]}
	}

	args := []string{
		"--config", fmt.Sprintf("%s/%s", configFilesDir, configName),
	}
	args = append(args, featureGatesArgs...)
	cmd := exec.Command(
		fmt.Sprintf("../%s", GetConfigVariable("OTEL_BINARY_FILE")),
		args...,
	)
	err := cmd.Start()
	require.NoError(t, err)

	t.Logf("Process started with PID: %d\n", cmd.Process.Pid)
	// wait for 5 seconds to allow the process to start
	time.Sleep(5 * time.Second)

	return cmd
}
