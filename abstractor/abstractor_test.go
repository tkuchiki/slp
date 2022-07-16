package abstractor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tkuchiki/slp/parser/sqlparser"
)

func TestSQLAbstractor_Abstract(t *testing.T) {
	abst := NewSQLAbstractor(false, false)
	parser := sqlparser.NewSQLParser()

	cases := []struct {
		original string
		want     string
	}{
		// SELECT
		{
			"SELECT * FROM `t` WHERE `c` = 1",
			"SELECT * FROM `t` WHERE `c` = N",
		},
		{
			"SELECT * FROM `t` WHERE `c` = 'abc'",
			"SELECT * FROM `t` WHERE `c` = 'S'",
		},
		{
			"SELECT * FROM `t` WHERE `c1` = 1 AND `c2` = 'abc'",
			"SELECT * FROM `t` WHERE `c1` = N AND `c2` = 'S'",
		},
		{
			"SELECT * FROM `t` WHERE `c1` = 1 OR `c2` = 'abc'",
			"SELECT * FROM `t` WHERE `c1` = N OR `c2` = 'S'",
		},
		{
			"SELECT * FROM `t` WHERE `c` IN ('abc')",
			"SELECT * FROM `t` WHERE `c` IN ('S')",
		},
		// INSERT
	}

	for _, tc := range cases {
		stmt, err := parser.Parse(tc.original)
		if err != nil {
			t.Errorf("Failed parsing: %v", err)
		}
		got, err := abst.Abstract(stmt)
		if err != nil {
			t.Errorf("Failed abstraction: %v", err)
		}

		if diff := cmp.Diff(got, tc.want); diff != "" {
			t.Errorf("(-got, +want)\n%s", diff)
		}
	}
}
