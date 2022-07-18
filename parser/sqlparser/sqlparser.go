package sqlparser

import (
	mparser "github.com/pingcap/tidb/parser"
	mast "github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
)

type SQLParser struct {
	p *mparser.Parser
}

func NewSQLParser() *SQLParser {
	return &SQLParser{
		p: mparser.New(),
	}
}

func (s *SQLParser) Parse(sql string) (*mast.StmtNode, error) {
	stmtNodes, _, err := s.p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	return &stmtNodes[0], nil
}
