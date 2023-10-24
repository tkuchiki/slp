package cmd

import (
	"strings"
	"testing"

	"github.com/tkuchiki/slp/internal/testutil"
)

func TestCommonFlags(t *testing.T) {
	tempDir := t.TempDir()

	slowlogFile := "../../../example/mysql.slow.log"

	tempConfig, err := testutil.CreateTempDirAndFile(tempDir, "test_common_flags_temp_config", testutil.ConfigFile())
	if err != nil {
		t.Fatal(err)
	}

	tempPos, err := testutil.CreateTempDirAndFile(tempDir, "test_common_flags_temp_pos", "")
	if err != nil {
		t.Fatal(err)
	}

	tempDump, err := testutil.CreateTempDirAndFile(tempDir, "test_common_flags_temp_dump", "")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		args []string
	}{
		{
			args: []string{"my", "--file", slowlogFile,
				"--noheaders",
				"--format", "tsv",
				"--config", tempConfig,
			},
		},
		{
			args: []string{"my", "--file", slowlogFile,
				"--filters", "Query matches 'SELECT'",
				"--format", "markdown",
				"--limit", "5",
				"--matching-groups", "SELECT .+",
				"--output", "count,query",
				"--page", "10",
				"--percentiles", "99",
				"--reverse",
				"--show-footers",
				"--sort", "query",
				"--bundle-where-in",
				"--bundle-values",
				"--noabstract",
			},
		},
		{
			args: []string{"my", "--file", slowlogFile,
				"-f", "Query matches 'SELECT'",
				"-m", "SELECT .+",
				"-o", "count,query",
				"-r",
				"-a",
			},
		},
		{
			args: []string{"my", "--file", slowlogFile,
				"--pos", tempPos,
			},
		},
		{
			args: []string{"my", "--file", slowlogFile,
				"--pos", tempPos,
				"--nosave-pos",
			},
		},
		// Do not change the order
		{
			args: []string{"my", "--file", slowlogFile,
				"--dump", tempDump,
			},
		},
		{
			args: []string{"my",
				"--load", tempDump,
			},
		},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			command := NewCommand("test")
			command.setArgs(tt.args)

			err := command.Execute()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
