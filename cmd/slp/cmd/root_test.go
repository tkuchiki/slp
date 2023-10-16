package cmd

import "testing"

func TestNewRootCmd(t *testing.T) {
	command := NewCommand("test")
	command.setArgs([]string{"--help"})

	err := command.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
