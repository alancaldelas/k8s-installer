package main

import (
	"os"
	"testing"
)

// Helper function to make a temp dir to run tests
func createYAML(t *testing.T, filename, content string) string {
	// Mark as helper function
	t.Helper()
	tempDir := t.TempDir()
	filePath := tempDir + "/" + filename

	err := os.WriteFile(filePath, []byte(content), 0644)

	if err != nil {
		t.Fatalf("failed to create temp file")
	}

	t.Log("Created temp file: ", filePath)

	return filePath

}

// Unit test for making sure we have a successful read of the config file
func TestReadConfig(t *testing.T) {
	data := `
Servers:
  - name: "Orange"
    port: 22
    Ip: "10.216.188.48"`

	fpath := createYAML(t, "test.yaml", data)

	t.Log("Reading Config file... ", fpath)
	// Run actual test
	_, err := ReadConfig(fpath)

	if err != nil {
		t.Errorf("failed to read the config file... %v", err)
	}
}
