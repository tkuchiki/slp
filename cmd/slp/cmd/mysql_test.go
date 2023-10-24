package cmd

import "testing"

func Test_newMySQLCmd(t *testing.T) {
	command := NewCommand("test")
	slowlogFile := "../../../example/mysql.slow.log"

	args := []string{"my",
		"--file", slowlogFile,
	}
	command.setArgs(args)

	err := command.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
