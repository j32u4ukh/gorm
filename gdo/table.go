package gdo

import (
	"github.com/j32u4ukh/gorm/stmt"
	"github.com/j32u4ukh/gorm/stmt/dialect"
)

type Table struct {
	*stmt.Table
	// 是否對 SQL injection 做處理
	useAntiInjection bool
}

func NewTable(name string, tableParam *stmt.TableParam, columnParams []*stmt.ColumnParam, engine string, collate string, dial dialect.SQLDialect) *Table {
	t := &Table{
		Table:            stmt.NewTable(name, tableParam, columnParams, engine, collate),
		useAntiInjection: false,
	}
	return t
}

func (t *Table) UseAntiInjection(active bool) {
	t.useAntiInjection = active
}
