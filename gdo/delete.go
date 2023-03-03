package gdo

import "github.com/pkg/errors"

func (t *Table) SetDeleteCondition(where *WhereStmt) {
	if t.useAntiInjection {
		where.UseAntiInjection()
	}
	t.DeleteStmt.SetCondition(where.ToStmtWhere())
}

func (t *Table) BuildDeleteStmt() (string, error) {
	sql, err := t.DeleteStmt.ToStmt()
	t.DeleteStmt.Release()

	if err != nil {
		return "", errors.Wrap(err, "生成 DeleteStmt 時發生錯誤")
	}

	return sql, nil
}
