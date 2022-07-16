package abstractor

import (
	"fmt"
	"strings"

	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/format"
	"github.com/tkuchiki/slp/visitor"
)

type SQLAbstractor struct {
	v            *visitor.Visitor
	restoreFlags format.RestoreFlags
}

func NewSQLAbstractor() *SQLAbstractor {
	return &SQLAbstractor{
		v:            visitor.NewVisitor(),
		restoreFlags: format.RestoreNameBackQuotes | format.RestoreSpacesAroundBinaryOperation,
	}
}

func (a *SQLAbstractor) Abstract(rootNode *ast.StmtNode) (string, error) {
	(*rootNode).Accept(a.v)

	var sb strings.Builder
	restoreCtx := format.NewRestoreCtx(a.restoreFlags, &sb)

	err := (*rootNode).Restore(restoreCtx)
	if err != nil {
		return "", err
	}
	fmt.Println(sb.String())

	return "", nil
}
