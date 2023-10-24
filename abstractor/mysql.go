package abstractor

import (
	"strings"

	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/format"
	"github.com/tkuchiki/slp/mysql/visitor"
)

type MySQLAbstractor struct {
	v            *visitor.Visitor
	restoreFlags format.RestoreFlags
}

func NewMySQLAbstractor(bundleWhereIn, bundleInsertValues bool) *MySQLAbstractor {
	return &MySQLAbstractor{
		v:            visitor.NewVisitor(bundleWhereIn, bundleInsertValues),
		restoreFlags: format.RestoreNameBackQuotes | format.RestoreSpacesAroundBinaryOperation,
	}
}

func (a *MySQLAbstractor) Abstract(rootNode *ast.StmtNode) (string, error) {
	(*rootNode).Accept(a.v)

	var sb strings.Builder
	restoreCtx := format.NewRestoreCtx(a.restoreFlags, &sb)

	err := (*rootNode).Restore(restoreCtx)
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}
