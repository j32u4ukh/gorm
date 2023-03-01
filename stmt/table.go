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
		for _, column := range t.CreateStmt.Columns {
			// if column.Default != "current_timestamp()" {
			t.ColumnNames.Append(column.Name)
			// }
		}
		t.InsertStmt.SetColumnNames(t.ColumnNames.Elements)
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

func (t *Table) String() string {
	info := fmt.Sprintf("Table %s", t.CreateStmt.TableName)

	for i, col := range t.Columns {
		info += fmt.Sprintf("\n%d) %s", i, col)
	}

	return info
}

// 根據 Sync 的需求，有需要再 Clone 即可
func (t *Table) SyncClone() *Table {
	clone := &Table{
		CreateStmt: t.CreateStmt.Clone(),
	}
	return clone
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Create
////////////////////////////////////////////////////////////////////////////////////////////////////
// 添加欄位
func (t *Table) AddColumn(column *Column) *Table {
	// 避免欄位重複添加
	if !t.ColumnNames.Contains(column.Name) {
		t.CreateStmt.AddColumn(column)

		if column.Default != "current_timestamp()" {
			t.ColumnNames.Append(column.Name)
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

func (t *Table) InsertWithCheck(datas map[string]string) error {
	err := t.checkInsertData(datas)

	if err != nil {
		return errors.Wrap(err, "欄位檢查時發現錯誤，中止數據插入")
	}

	t.Insert(datas)
	return nil
}

func (t *Table) checkInsertData(datas map[string]string) error {
	var value string

	// 移除不屬於當前表格的欄位名稱
	for value = range datas {
		if !t.ColumnNames.Contains(value) {
			delete(datas, value)
		}
	}

	if len(datas) == 0 {
		return errors.New("傳入數據皆不屬於當前表格的欄位名稱")
	}

	var ok bool
	var column *Column

	// 檢查當前表格的欄位數據有效性
	for _, column = range t.Columns {
		if value, ok = datas[column.Name]; !ok {
			// datas 內沒有包含此欄位的數據

			// 若此欄位是 PrimaryKey，且非自動增加類型
			if column.IsPrimaryKey && column.Default != "AI" {
				return errors.New(fmt.Sprintf("缺少 PrimaryKey 欄位(%s)", column.Name))
			}

			// 此欄位填入預設值
			value = column.Default

		} else {
			// datas 內有此欄位的數據

			// ===== 檢查數據是否有效 =====
			if !(column.Default == "AI" || column.Default == "current_timestamp()") {
				if value == "NULL" {
					value = column.Default
				}
			}
			// ===========================
		}

		// 將修正完的數值，再放入 datas 中
		datas[column.Name] = value
	}

	return nil
}

func (t *Table) Insert(datas map[string]string) error {
	// 確保 insertStmt 有語法生成用的欄位名稱
	if t.InsertStmt.ColumnStmt == "" {
		t.InsertStmt.SetColumnNames(t.ColumnNames.Elements)
	}

	insertData := []string{}
	var column *Column

	for _, column = range t.Columns {
		insertData = append(insertData, datas[column.Name])
	}

	t.InsertStmt.Insert(insertData)
	return nil
}

func (t *Table) InsertRawData(datas []string) error {
	// 確保 insertStmt 有語法生成用的欄位名稱
	if t.InsertStmt.ColumnStmt == "" {
		t.InsertStmt.SetColumnNames(t.ColumnNames.Elements)
	}
	t.InsertStmt.Insert(datas)
	return nil
}

// 取得緩存數量
func (t *Table) GetInsertBufferNumber() int32 {
	return t.InsertStmt.GetBufferNumber()
}

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
	var col *Column

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

////////////////////////////////////////////////////////////////////////////////////////////////////
// Column
////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *Table) GetColumnNumber() int32 {
	return int32(len(t.Columns))
}

func (t *Table) SetColumn(idx int32, column *Column) error {
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

func (t *Table) GetColumn(idx int32) *Column {
	if idx < int32(len(t.Columns)) {
		return t.Columns[idx]
	}
	return nil
}

func (t *Table) GetColumnByName(name string) *Column {
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
func (t *Table) GetPrimaryColumn() (*Column, error) {
	for _, col := range t.Columns {
		if col.IsPrimaryKey {
			return col, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Not found primary key in %s", t.CreateStmt.TableName))
}
