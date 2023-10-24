package abstractor

import (
	"strings"

	pgquery "github.com/pganalyze/pg_query_go/v4"
)

type PGAbstractor struct {
	noAbstract    bool
	bundleWhereIn bool
	bundleValues  bool
}

func NewPGAbstractor(noAbstract, bundleWhereIn, bundleValues bool) *PGAbstractor {
	return &PGAbstractor{
		noAbstract:    noAbstract,
		bundleWhereIn: bundleWhereIn,
		bundleValues:  bundleValues,
	}
}

func (a *PGAbstractor) Abstract(sql string) (string, error) {
	if a.noAbstract && !a.bundleWhereIn && !a.bundleValues {
		return sql, nil
	}

	result, err := pgquery.Parse(sql)
	if err != nil {
		return "", err
	}

	for _, stmt := range result.Stmts {
		a.parseStmt(stmt.GetStmt())
	}

	output, err := pgquery.Deparse(result)
	if err != nil {
		return "", err
	}

	if !a.noAbstract {
		output = strings.ReplaceAll(output, "'N'", "N")
	}

	return output, nil
}

func (a *PGAbstractor) parseStmt(stmt *pgquery.Node) {
	if stmt == nil {
		return
	}

	switch _stmt := stmt.Node.(type) {
	case *pgquery.Node_SelectStmt:
		a.parseSelectStmt(_stmt.SelectStmt)
	case *pgquery.Node_InsertStmt:
		a.parseInsertStmt(_stmt.InsertStmt)
	case *pgquery.Node_UpdateStmt:
		a.parseUpdateStmt(_stmt.UpdateStmt)
	case *pgquery.Node_DeleteStmt:
		a.parseDeleteStmt(_stmt.DeleteStmt)
	case *pgquery.Node_CopyStmt:
		a.parseCopyStmt(_stmt.CopyStmt)
	}
}

func (a *PGAbstractor) parseSelectStmt(stmt *pgquery.SelectStmt) {
	if stmt == nil {
		return
	}

	a.parseNode(stmt.GetWhereClause())
	a.parseWithClause(stmt.GetWithClause())

	if stmt.ValuesLists != nil {
		a.parseNodes(stmt.GetValuesLists())

		if a.bundleValues {
			values := stmt.GetValuesLists()
			stmt.ValuesLists = values[:1]
		}
	}
	a.parseNode(stmt.GetLimitCount())
}

func (a *PGAbstractor) parseInsertStmt(stmt *pgquery.InsertStmt) {
	if stmt == nil {
		return
	}

	a.parseNode(stmt.GetSelectStmt())
	if stmt.SelectStmt != nil {
		a.parseWithClause(stmt.SelectStmt.GetWithClause())
	}

	a.parseWithClause(stmt.GetWithClause())
	a.parseOnConflictClause(stmt.GetOnConflictClause())
}

func (a *PGAbstractor) parseUpdateStmt(stmt *pgquery.UpdateStmt) {
	if stmt == nil {
		return
	}

	a.parseNode(stmt.GetWhereClause())
	a.parseNodes(stmt.GetTargetList())
	a.parseWithClause(stmt.GetWithClause())
}

func (a *PGAbstractor) parseDeleteStmt(stmt *pgquery.DeleteStmt) {
	if stmt == nil {
		return
	}

	a.parseNode(stmt.GetWhereClause())
	a.parseWithClause(stmt.GetWithClause())
	a.parseNodes(stmt.GetUsingClause())
}

func (a *PGAbstractor) parseCopyStmt(stmt *pgquery.CopyStmt) {
	if stmt == nil {
		return
	}

	a.parseNode(stmt.GetQuery())
}

func (a *PGAbstractor) parseNode(node *pgquery.Node) {
	if node == nil {
		return
	}

	switch _node := node.Node.(type) {
	case *pgquery.Node_BoolExpr:
		a.parseNodes(_node.BoolExpr.Args)
	case *pgquery.Node_AExpr:
		a.parseAExpr(_node.AExpr)
	case *pgquery.Node_CommonTableExpr:
		a.parseNode(_node.CommonTableExpr.GetCtequery())
	case *pgquery.Node_AArrayExpr:
		a.parseNodes(_node.AArrayExpr.GetElements())
	case *pgquery.Node_AConst:
		a.parseAConst(_node.AConst)
	case *pgquery.Node_WithClause:
		a.parseNodes(_node.WithClause.Ctes)
	case *pgquery.Node_SelectStmt:
		a.parseSelectStmt(_node.SelectStmt)
	case *pgquery.Node_InsertStmt:
		a.parseInsertStmt(_node.InsertStmt)
	case *pgquery.Node_UpdateStmt:
		a.parseUpdateStmt(_node.UpdateStmt)
	case *pgquery.Node_DeleteStmt:
		a.parseDeleteStmt(_node.DeleteStmt)
	case *pgquery.Node_CopyStmt:
		a.parseCopyStmt(_node.CopyStmt)
	case *pgquery.Node_List:
		a.parseNodes(_node.List.GetItems())
	case *pgquery.Node_ResTarget:
		a.parseNode(_node.ResTarget.GetVal())
	case *pgquery.Node_SubLink:
		a.parseNode(_node.SubLink.GetSubselect())
	case *pgquery.Node_RangeSubselect:
		a.parseNode(_node.RangeSubselect.GetSubquery())
	}
}

func (a *PGAbstractor) parseAExpr(aExpr *pgquery.A_Expr) {
	if aExpr.Lexpr != nil {
		switch aExpr.Lexpr.Node.(type) {
		case *pgquery.Node_AConst:
			a.parseNode(aExpr.Lexpr)
		}
	}

	if aExpr.Rexpr != nil {
		a.parseNode(aExpr.Rexpr)

		if a.bundleWhereIn {
			switch node := aExpr.Rexpr.Node.(type) {
			case *pgquery.Node_List:
				items := node.List.GetItems()
				node.List.Items = items[:1]
			}
		}
	}
}

func (a *PGAbstractor) parseAConst(aConst *pgquery.A_Const) {
	if aConst == nil || a.noAbstract {
		return
	}

	switch aConst.Val.(type) {
	case *pgquery.A_Const_Ival, *pgquery.A_Const_Fval:
		aConst.Val = &pgquery.A_Const_Sval{
			Sval: &pgquery.String{Sval: "N"},
		}
	case *pgquery.A_Const_Sval:
		aConst.Val = &pgquery.A_Const_Sval{
			Sval: &pgquery.String{Sval: "S"},
		}
	}
}

func (a *PGAbstractor) parseNodes(nodes []*pgquery.Node) {
	if nodes == nil {
		return
	}

	for _, node := range nodes {
		a.parseNode(node)
	}
}

func (a *PGAbstractor) parseWithClause(with *pgquery.WithClause) {
	if with == nil {
		return
	}

	for _, node := range with.GetCtes() {
		a.parseNode(node)
	}
}

func (a *PGAbstractor) parseOnConflictClause(onConflict *pgquery.OnConflictClause) {
	if onConflict == nil {
		return
	}

	a.parseNode(onConflict.GetWhereClause())
	a.parseNodes(onConflict.GetTargetList())
}
