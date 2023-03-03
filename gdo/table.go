package gdo

import (
	"fmt"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gorm/stmt"
	"github.com/j32u4ukh/gorm/stmt/dialect"
)

type Table struct {
	*stmt.CreateStmt
	*stmt.InsertStmt
	*stmt.SelectStmt
	*stmt.UpdateStmt
	*stmt.DeleteStmt
	//////////////////////////////////////////////////
	// 是否對 SQL injection 做處理
	useAntiInjection bool
	ColumnNames      *cntr.Array[string]
	nColumn          int32
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Table
////////////////////////////////////////////////////////////////////////////////////////////////////
func NewTable(name string, tableParam *stmt.TableParam, columnParams []*stmt.ColumnParam, engine string, collate string, dial dialect.SQLDialect) *Table {
	t := &Table{
		CreateStmt:       stmt.NewCreateStmt(name, tableParam, columnParams, engine, collate),
		InsertStmt:       stmt.NewInsertStmt(name),
		SelectStmt:       stmt.NewSelectStmt(name),
		UpdateStmt:       stmt.NewUpdateStmt(name),
		DeleteStmt:       stmt.NewDeleteStmt(name),
		useAntiInjection: false,
		ColumnNames:      cntr.NewArray[string](),
	}
	if len(t.CreateStmt.Columns) > 0 {
		// 會自行賦值的欄位也需填入 NULL，因此所有欄位名稱都要求填入
		for _, column := range t.CreateStmt.Columns {
			if column.IgnoreThis {
				continue
			}
			t.ColumnNames.Append(column.Name)
		}
		t.nColumn = int32(t.ColumnNames.Length())
		t.InsertStmt.SetColumnNames(t.ColumnNames.Elements)
		fmt.Printf("gdo.NewTable | ColumnNames(%d): %+v\n", t.nColumn, t.ColumnNames.Elements)
	}
	return t
}

func (t *Table) SetDbName(dbName string) *Table {
	t.CreateStmt.SetDbName(dbName)
	t.InsertStmt.SetDbName(dbName)
	t.SelectStmt.SetDbName(dbName)
	t.UpdateStmt.SetDbName(dbName)
	t.DeleteStmt.SetDbName(dbName)
	return t
}

func (t *Table) GetDbName() string {
	return t.CreateStmt.DbName
}

func (t *Table) GetTableName() string {
	return t.CreateStmt.TableName
}

func (t *Table) UseAntiInjection(active bool) {
	t.useAntiInjection = active
}

func (t *Table) String() string {
	info := fmt.Sprintf("Table %s", t.CreateStmt.TableName)

	for i, col := range t.Columns {
		info += fmt.Sprintf("\n%d) %s", i, col)
	}

	return info
}

// NOTE: 根據 Sync 的需求，有需要再 Clone 即可
func (t *Table) SyncClone() *Table {
	clone := &Table{
		CreateStmt: t.CreateStmt.Clone(),
	}
	return clone
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Create
////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *Table) BuildCreateStmt() (string, error) {
	return t.CreateStmt.ToStmt()
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Sync
////////////////////////////////////////////////////////////////////////////////////////////////////
// 根據傳入的欄位名稱列表 orders，對欄位進行重新排序
// return
// 	- 被修改順序的欄位名稱
func (t *Table) RefreshColumnOrder(orders []string) *cntr.Array[string] {
	changes := cntr.NewArray[string]()
	changed := t.refreshColumnOrder(orders)

	for changed != "" {
		changes.Append(changed)
		changed = t.refreshColumnOrder(orders)
	}
	return changes
}

func (t *Table) refreshColumnOrder(orders []string) string {
	changed := ""
	var i, j int
	var order string = ""
	var col *stmt.Column

	for i, order = range orders {
		for j, col = range t.Columns {
			if col.Name == order {
				if i != j {
					changed = order
				}
				break
			}
		}
		if changed != "" {
			break
		}
	}
	if changed != "" {
		col = t.Columns[i]
		t.Columns[i] = t.Columns[j]
		t.Columns[j] = col
	}
	return changed
}
