package visitor

import (
	"github.com/k0kubun/pp"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
)

type Visitor struct {
}

func NewVisitor() *Visitor {
	return &Visitor{}
}

func (v *Visitor) Enter(in ast.Node) (ast.Node, bool) {
	switch stmt := in.(type) {
	case *ast.SelectStmt:
		parseSelectStmt(stmt)
	case *ast.InsertStmt:
		pp.Println(stmt)
		parseInsertStmt(stmt)
	case *ast.UpdateStmt:
		parseUpdateStmt(stmt)
	case *ast.DeleteStmt:
		parseDeleteStmt(stmt)
	}
	return in, false
}

func (v *Visitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func parseSelectStmt(in *ast.SelectStmt) {
	if in.Where != nil {
		parseExpr(in.Where)
	}

	if in.Having != nil {
		parseExpr(in.Having.Expr)
	}

	if in.Limit != nil {
		parseAndAbstractLimit(in.Limit)
	}
}

func parseExpr(in ast.ExprNode) {
	switch expr := in.(type) {
	case *ast.BinaryOperationExpr:
		switch l := expr.L.(type) {
		case *ast.BinaryOperationExpr:
			parseExpr(l)
		}

		switch r := expr.R.(type) {
		case *test_driver.ValueExpr:
			abstractValue(r)
		case *ast.BinaryOperationExpr:
			parseExpr(r)
		}
	case *test_driver.ValueExpr:
		abstractValue(expr)
	}
}

func abstractValue(val *test_driver.ValueExpr) {
	kind := val.Datum.Kind()

	switch kind {
	case test_driver.KindInt64, test_driver.KindUint64, test_driver.KindFloat32, test_driver.KindFloat64, test_driver.KindMysqlDecimal:
		val.Datum = test_driver.NewStringDatum("N")
	default:
		if test_driver.KindNull != kind {
			val.Datum = test_driver.NewStringDatum("'S'")
		}
	}
	val.Type.SetCharset("")
}

func parseAndAbstractLimit(in *ast.Limit) {
	in.Count.(*test_driver.ValueExpr).Datum = test_driver.NewStringDatum("N")
	in.Count.(*test_driver.ValueExpr).Type.SetCharset("")
}

func parseInsertStmt(in *ast.InsertStmt) {
	if in.Lists != nil {
		for _, values := range in.Lists {
			for _, val := range values {
				parseExpr(val)
			}
		}
	}

	if in.OnDuplicate != nil {
		for _, assignment := range in.OnDuplicate {
			parseExpr(assignment.Expr)
		}
	}
}

func parseUpdateStmt(in *ast.UpdateStmt) {
	if in.List != nil {
		for _, assignment := range in.List {
			parseExpr(assignment.Expr)
		}
	}

	if in.Where != nil {
		parseExpr(in.Where)
	}

	if in.Limit != nil {
		parseAndAbstractLimit(in.Limit)
	}
}

func parseDeleteStmt(in *ast.DeleteStmt) {
	if in.Where != nil {
		parseExpr(in.Where)
	}

	if in.Limit != nil {
		parseAndAbstractLimit(in.Limit)
	}
}
