package cmd

import (
	"io"
	"os"
	"testing"
)

func executeForOutput(t *testing.T, args []string) string {
	t.Helper()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	originalStdout := os.Stdout
	os.Stdout = writer

	command := NewCommand("test")
	command.setArgs(args)
	executeErr := command.Execute()
	closeErr := writer.Close()
	os.Stdout = originalStdout

	output, readErr := io.ReadAll(reader)
	_ = reader.Close()
	if executeErr != nil {
		t.Fatal(executeErr)
	}
	if closeErr != nil {
		t.Fatal(closeErr)
	}
	if readErr != nil {
		t.Fatal(readErr)
	}

	return string(output)
}

func TestDumpAndLoadPreserveProfileOutput(t *testing.T) {
	dumpFile := t.TempDir() + "/profile.yaml"
	outputOptions := []string{
		"--format", "csv",
		"--output", "count,query,min-query-time,max-query-time,avg-lock-time,avg-rows-examined",
	}
	profileArgs := append([]string{
		"my",
		"--file", "../../../example/mysql.slow.log",
		"--dump", dumpFile,
	}, outputOptions...)
	loadArgs := append([]string{
		"my",
		"--load", dumpFile,
	}, outputOptions...)

	profileOutput := executeForOutput(t, profileArgs)
	loadOutput := executeForOutput(t, loadArgs)
	if loadOutput != profileOutput {
		t.Fatalf("loaded output:\n%s\nwant profiled output:\n%s", loadOutput, profileOutput)
	}
}
