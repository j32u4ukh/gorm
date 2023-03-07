package stmt

import (
	"fmt"

	"github.com/j32u4ukh/cntr"
	"github.com/pkg/errors"
)

type Table struct {
	*CreateStmt
	*InsertStmt
	*SelectStmt
	*UpdateStmt
	*DeleteStmt
	ColumnNames *cntr.Array[string]
	nColumn     int32
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Table
////////////////////////////////////////////////////////////////////////////////////////////////////
func NewTable(name string, tableParam *TableParam, columnParams []*ColumnParam, engine string, collate string) *Table {
	t := &Table{
		CreateStmt:  NewCreateStmt(name, tableParam, columnParams, engine, collate),
		InsertStmt:  NewInsertStmt(name),
		SelectStmt:  NewSelectStmt(name),
		UpdateStmt:  NewUpdateStmt(name),
		DeleteStmt:  NewDeleteStmt(name),
		ColumnNames: cntr.NewArray[string](),
	}
	if len(t.CreateStmt.Columns) > 0 {
		// 會自行賦值的欄位也需填入 NULL，因此所有欄位名稱都要求填入
		for _, column := range t.CreateStmt.Columns {
			if column.IgnoreThis {
				continue
			}
			t.ColumnNames.Append(column.Name)
		}

	}
	return t
}

func (t *Table) SetDbName(dbName string) {
	t.CreateStmt.SetDbName(dbName)
	t.InsertStmt.SetDbName(dbName)
	t.SelectStmt.SetDbName(dbName)
	t.UpdateStmt.SetDbName(dbName)
	t.DeleteStmt.SetDbName(dbName)
}

func (t *Table) GetDbName() string {
	return t.CreateStmt.DbName
}

func (t *Table) GetTableName() string {
	return t.CreateStmt.TableName
}

func (t *Table) String() string {
	info := fmt.Sprintf("Table %s", t.CreateStmt.TableName)

	for i, col := range t.Columns {
		info += fmt.Sprintf("\n%d) %s", i, col)
	}

	return info
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Create
////////////////////////////////////////////////////////////////////////////////////////////////////
// 添加欄位
func (t *Table) AddColumn(column *Column) *Table {
	// 避免欄位重複添加
	if !column.IgnoreThis && !t.ColumnNames.Contains(column.Name) {
		t.CreateStmt.AddColumn(column)

		if column.Default != "current_timestamp()" {
			t.ColumnNames.Append(column.Name)
			t.nColumn = int32(t.ColumnNames.Length())
			t.InsertStmt.SetColumnNames(t.ColumnNames.Elements)
		}
	}
	return t
}

func (t *Table) BuildCreateStmt() (string, error) {
	return t.CreateStmt.ToStmt()
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Insert
////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *Table) BuildInsertStmt() (string, error) {
	sql, err := t.InsertStmt.ToStmt()
	t.InsertStmt.Release()

	if err != nil {
		return "", errors.Wrap(err, "生成 InsertStmt 時發生錯誤")
	}

	return sql, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Select
////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *Table) SetSelectCondition(where *WhereStmt) {
	t.SelectStmt.SetCondition(where)
}

func (t *Table) BuildSelectStmt() (string, error) {
	sql, err := t.SelectStmt.ToStmt()
	t.SelectStmt.Release()
	if err != nil {
		return "", errors.Wrap(err, "生成 SelectStmt 時發生錯誤")
	}
	return sql, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Update
////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *Table) SetUpdateCondition(where *WhereStmt) {
	t.UpdateStmt.SetCondition(where)
}

func (t *Table) BuildUpdateStmt() (string, error) {
	sql, err := t.UpdateStmt.ToStmt()
	t.UpdateStmt.Release()

	if err != nil {
		return "", errors.Wrap(err, "生成 UpdateStmt 時發生錯誤")
	}

	return sql, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Delete
////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *Table) SetDeleteCondition(where *WhereStmt) {
	t.DeleteStmt.SetCondition(where)
}

func (t *Table) BuildDeleteStmt() (string, error) {
	sql, err := t.DeleteStmt.ToStmt()
	t.DeleteStmt.Release()

	if err != nil {
		return "", errors.Wrap(err, "生成 DeleteStmt 時發生錯誤")
	}

	return sql, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Column
////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *Table) GetColumnNumber() int32 {
	return t.nColumn
}

func (t *Table) SetColumn(idx int32, column *Column) error {
	if t.nColumn-1 < idx {
		return errors.New(fmt.Sprintf("idx(%d) out of length(%d).", idx, t.nColumn-1))
	}

	origin := t.Columns[idx]

	// 更新 tableParam 當中的欄位名稱(origin.Name -> column.Name)
	if origin.Name != column.Name {
		t.TableParam.UpdateIndexName(origin.Name, column.Name)
	}

	t.Columns[idx] = column
	return nil
}

func (t *Table) GetColumn(idx int32) *Column {
	if idx < int32(len(t.Columns)) {
		return t.Columns[idx]
	}
	return nil
}

func (t *Table) GetColumnByName(name string) (int32, *Column) {
	for i, column := range t.Columns {
		if column.Name == name {
			return int32(i), column
		}
	}
	return -1, nil
}

func (t *Table) GetColumnNames() *cntr.Array[string] {
	return t.ColumnNames.Clone()
}

func (t *Table) GetPrimaryColumn() []*Column {
	pks := []*Column{}
	for _, col := range t.Columns {
		if col.IsPrimaryKey {
			pks = append(pks, col)
		}
	}
	return pks
}
