package visitor

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
)

type Visitor struct {
	assembleWhereIn      bool
	assembleInsertValues bool
}

func NewVisitor(assembleWhereIn, assembleInsertValues bool) *Visitor {
	return &Visitor{
		assembleWhereIn:      assembleWhereIn,
		assembleInsertValues: assembleInsertValues,
	}
}

func (v *Visitor) Enter(in ast.Node) (ast.Node, bool) {
	switch stmt := in.(type) {
	case *ast.InsertStmt:
		parseInsertStmt(stmt, v.assembleInsertValues)
	case *test_driver.ValueExpr:
		abstractValue(stmt)
	case *ast.PatternInExpr:
		for _, val := range stmt.List {
			parseExpr(val)

			if v.assembleWhereIn {
				stmt.List = stmt.List[:1]
				break
			}
		}
	}
	return in, false
}

func (v *Visitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func parseExpr(in ast.ExprNode) {
	switch expr := in.(type) {
	case *ast.BinaryOperationExpr:
		switch l := expr.L.(type) {
		case *ast.BinaryOperationExpr:
			parseExpr(l)
		case *test_driver.ValueExpr:
			abstractValue(l)
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

func parseInsertStmt(in *ast.InsertStmt, assemble bool) {
	if !assemble {
		return
	}

	if in.Lists != nil {
		for _, values := range in.Lists {
			for _, val := range values {
				parseExpr(val)
			}
			if assemble {
				break
			}
		}

		if assemble {
			in.Lists = in.Lists[:1]
		}
	}

	if in.OnDuplicate != nil {
		for _, assignment := range in.OnDuplicate {
			parseExpr(assignment.Expr)
		}
	}
}
