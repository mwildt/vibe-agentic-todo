package tests

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

// TestHelpCommand tests that CLI commands support help options
func TestHelpCommand(t *testing.T) {
	// Test root command help
	rootHelpCmd := exec.Command("go", "run", "vibe-agentic/cmd/cli", "--help")
	rootHelpCmd.Env = append(os.Environ(), "GO111MODULE=on")
	rootHelpCmd.Dir = ".."
	var rootStdout, rootStderr bytes.Buffer
	rootHelpCmd.Stdout = &rootStdout
	rootHelpCmd.Stderr = &rootStderr
	
	if err := rootHelpCmd.Run(); err != nil {
		t.Fatalf("root help command failed: %v, stderr: %s", err, rootStderr.String())
	}
	
	// Check that help output contains expected information
	if !bytes.Contains(rootStdout.Bytes(), []byte("A command-line interface for managing the Vibe application")) {
		t.Errorf("root help does not contain expected description: got %s", rootStdout.String())
	}
	
	// Test user command help
	userHelpCmd := exec.Command("go", "run", "vibe-agentic/cmd/cli", "user", "--help")
	userHelpCmd.Env = append(os.Environ(), "GO111MODULE=on")
	userHelpCmd.Dir = ".."
	var userStdout, userStderr bytes.Buffer
	userHelpCmd.Stdout = &userStdout
	userHelpCmd.Stderr = &userStderr
	
	if err := userHelpCmd.Run(); err != nil {
		t.Fatalf("user help command failed: %v, stderr: %s", err, userStderr.String())
	}
	
	// Check that help output contains expected information
	if !bytes.Contains(userStdout.Bytes(), []byte("Commands for managing users in the YAML configuration")) {
		t.Errorf("user help does not contain expected description: got %s", userStdout.String())
	}
	
	// Test user add command help
	addHelpCmd := exec.Command("go", "run", "vibe-agentic/cmd/cli", "user", "add", "--help")
	addHelpCmd.Env = append(os.Environ(), "GO111MODULE=on")
	addHelpCmd.Dir = ".."
	var addStdout, addStderr bytes.Buffer
	addHelpCmd.Stdout = &addStdout
	addHelpCmd.Stderr = &addStderr
	
	if err := addHelpCmd.Run(); err != nil {
		t.Fatalf("user add help command failed: %v, stderr: %s", err, addStderr.String())
	}
	
	// Check that help output contains expected information
	if !bytes.Contains(addStdout.Bytes(), []byte("Add a new user to the YAML configuration file")) {
		t.Errorf("user add help does not contain expected description: got %s", addStdout.String())
	}
	
	// Check that help output contains flag information
	if !bytes.Contains(addStdout.Bytes(), []byte("--username")) || !bytes.Contains(addStdout.Bytes(), []byte("--password")) {
		t.Errorf("user add help does not contain expected flags: got %s", addStdout.String())
	}
}
