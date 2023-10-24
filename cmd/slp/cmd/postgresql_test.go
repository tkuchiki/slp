package cmd

import "testing"

func Test_newPostgreSQLCmd(t *testing.T) {
	command := NewCommand("test")
	slowlogFile := "../../../example/postgresql.slow.log"

	args := []string{"pg",
		"--file", slowlogFile,
		"--log-line-prefix", "%m",
	}
	command.setArgs(args)

	err := command.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
