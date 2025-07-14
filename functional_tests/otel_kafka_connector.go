package functional_tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func stopOTelKafkaConnector(t *testing.T, cmd *exec.Cmd) {
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
		fmt.Println("Process force killed after timeout.")
	case err := <-done:
		require.NoError(t, err)
		fmt.Println("Process exited gracefully.")
	}
}

func startOTelKafkaConnector(t *testing.T, configName string, configFilesDir string) *exec.Cmd {
	// Start the process
	cmd := exec.Command(fmt.Sprintf("./%s", GetConfigVariable("OTEL_BINARY_FILE")), "--config", fmt.Sprintf("%s/%s", configFilesDir, configName))

	err := cmd.Start()
	require.NoError(t, err)

	fmt.Printf("Process started with PID: %d\n", cmd.Process.Pid)
	// wait for 5 seconds to allow the process to start
	time.Sleep(5 * time.Second)

	return cmd
}
