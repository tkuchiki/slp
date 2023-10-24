package sqlparser

import (
	mparser "github.com/pingcap/tidb/parser"
	mast "github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
)

type MySQLSQLParser struct {
	p *mparser.Parser
}

func NewMySQLSQLParser() *MySQLSQLParser {
	return &MySQLSQLParser{
		p: mparser.New(),
	}
}

func (s *MySQLSQLParser) Parse(sql string) (*mast.StmtNode, error) {
	stmtNodes, _, err := s.p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	return &stmtNodes[0], nil
}
