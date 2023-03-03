package gdo

import (
	"reflect"

	"github.com/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// Update
////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *Table) Update(key string, value any, ptrToDb func(reflect.Value, bool) string) {
	t.UpdateStmt.Update(key, ValueToDb(reflect.ValueOf(value), t.useAntiInjection, ptrToDb))
}

func (t *Table) UpdateRawData(key string, value string) {
	t.UpdateStmt.Update(key, value)
}

func (t *Table) SetUpdateCondition(where *WhereStmt) {
	if t.useAntiInjection {
		where.UseAntiInjection()
	}
	t.UpdateStmt.SetCondition(where.ToStmtWhere())
}

func (t *Table) BuildUpdateStmt() (string, error) {
	sql, err := t.UpdateStmt.ToStmt()
	t.UpdateStmt.Release()

	if err != nil {
		return "", errors.Wrap(err, "生成 UpdateStmt 時發生錯誤")
	}

	return sql, nil
}
