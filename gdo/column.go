package gdo

import (
	"fmt"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gorm/stmt"

	"github.com/pkg/errors"
)

// 添加欄位
func (t *Table) AddColumn(column *stmt.Column) *Table {
	// 避免欄位重複添加
	if !t.ColumnNames.Contains(column.Name) {
		t.CreateStmt.AddColumn(column)
		t.ColumnNames.Append(column.Name)
	}
	return t
}

func (t *Table) GetColumnNumber() int32 {
	return int32(len(t.Columns))
}

func (t *Table) SetColumn(idx int32, column *stmt.Column) error {
	if t.GetColumnNumber()-1 < idx {
		return errors.New(fmt.Sprintf("idx(%d) out of length(%d).", idx, t.GetColumnNumber()-1))
	}

	origin := t.Columns[idx]

	// 更新 tableParam 當中的欄位名稱(origin.Name -> column.Name)
	if origin.Name != column.Name {
		t.TableParam.UpdateIndexName(origin.Name, column.Name)
	}

	t.Columns[idx] = column
	return nil
}

func (t *Table) GetColumn(idx int32) *stmt.Column {
	if idx < int32(len(t.Columns)) {
		return t.Columns[idx]
	}
	return nil
}

func (t *Table) GetColumnByName(name string) *stmt.Column {
	for _, column := range t.Columns {
		if column.Name == name {
			return column
		}
	}
	return nil
}

func (t *Table) GetIndexByName(name string) int32 {
	for i, column := range t.Columns {
		if column.Name == name {
			return int32(i)
		}
	}
	return -1
}

func (t *Table) GetColumnNames() *cntr.Array[string] {
	arr := cntr.NewArray[string]()

	for _, column := range t.Columns {
		arr.Append(column.Name)
	}

	return arr
}

// TODO: PrimaryKey 可以有多組
func (t *Table) GetPrimaryColumn() (*stmt.Column, error) {
	for _, col := range t.Columns {
		if col.IsPrimaryKey {
			return col, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Not found primary key in %s", t.CreateStmt.TableName))
}
