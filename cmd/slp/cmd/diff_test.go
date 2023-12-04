package cmd

import "testing"

func Test_newDiffCmd(t *testing.T) {
	command := NewCommand("test")
	slpDumpFile1 := "../../../example/slp-dump-1.yaml"
	slpDumpFile2 := "../../../example/slp-dump-2.yaml"

	args := []string{"diff",
		slpDumpFile1,
		slpDumpFile2,
		"--show-footers",
	}
	command.setArgs(args)

	err := command.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
