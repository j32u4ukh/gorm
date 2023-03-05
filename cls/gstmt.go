package cls

import (
	"fmt"

	"github.com/j32u4ukh/gorm/gdo"
	"github.com/j32u4ukh/gorm/stmt/dialect"
	"github.com/pkg/errors"
)

const (
	TIME_LAYOUT string = "2006-01-02 15:04:05"
)

var gstmtMap map[byte]*Gstmt

type Gstmt struct {
	DbName string
	Dial   dialect.SQLDialect
	tables map[byte]*StructTable
}

func init() {
	gstmtMap = map[byte]*Gstmt{}
}

func SetGstmt(index byte, dbName string, dialect dialect.SQLDialect) (*Gstmt, error) {
	if _, ok := gstmtMap[index]; ok {
		return nil, errors.New(fmt.Sprintf("index: %d 已使用", index))
	}
	g, err := NewGstmt(dbName, dialect)
	if err != nil {
		return nil, errors.Wrap(err, "建立 Gstmt 時發生錯誤")
	}
	gstmtMap[index] = g
	return g, nil
}

func GetGstmt(index byte) (*Gstmt, error) {
	if g, ok := gstmtMap[index]; ok {
		return g, nil
	}
	return nil, errors.New(fmt.Sprintf("index: %d 沒有對應的 *Gstmt", index))
}

// 由外部定義實際的 DB 實體，GoSql 則對該 DB 實體進行操作
// db: DB 指標
// dialect: 資料庫語言(mysql/maria/...)
func NewGstmt(dbName string, dial dialect.SQLDialect) (*Gstmt, error) {
	g := &Gstmt{
		DbName: dbName,
		Dial:   dial,
		tables: map[byte]*StructTable{},
	}
	return g, nil
}

func (g *Gstmt) SetDbName(dbName string) {
	g.DbName = dbName
}

func (g *Gstmt) UseAntiInjection(active bool) {
	for _, t := range g.tables {
		t.UseAntiInjection(active)
	}
}

// ====================================================================================================
// Create
// ====================================================================================================
// 於 Gstmt 中建立表格結構
// - tid: 表格 ID
// - folder: 資料夾位置
// - tableName: 表格名稱
// return
// - 生成表格之 SQL 語法
// - 錯誤訊息
func (g *Gstmt) CreateTable(tid byte, tableName string, dial dialect.SQLDialect, data any) (string, error) {
	// 檢查表格是否已存在變數中
	if _, ok := g.tables[tid]; ok {
		fmt.Printf("(g *Gstmt) CreateTable | 表格 %s 已存在\n", tableName)
		return "", nil
	}

	// 根據檔案名稱，取得表格與欄位參數
	tableParam, colParams, err := GetParams(data, dial)

	if err != nil {
		return "", errors.Wrapf(err, "從 %s.proro 讀取參數時發生錯誤", tableName)
	}

	// fmt.Printf("(g *Gstmt) CreateTable | tableParam: %+v.\n", tableParam)
	t := NewStructTable(tableName, tableParam, colParams, g.Dial)
	fmt.Printf("(g *Gstmt) CreateTable | table: %+v.\n", t)
	t.SetDbName(g.DbName)

	// 將 Table 加入管理，可利用 tableName 進行存取
	g.tables[tid] = t
	sql, _ := t.BuildCreateStmt()
	return sql, nil
}

// ====================================================================================================
// Insert
// ====================================================================================================
func (g *Gstmt) Insert(tid byte, datas []any) (string, error) {
	if table, ok := g.tables[tid]; ok {
		table.InitByStruct(datas[0])
		table.Insert(datas)
		sql, err := table.BuildInsertStmt()
		if err != nil {
			return "", errors.Wrap(err, "Failed to build InsertStmt.")
		}
		return sql, nil
	}
	return "", errors.New(fmt.Sprintf("找不到編號為 %d 的表格", tid))
}

// ====================================================================================================
// Select
// ====================================================================================================

func (g *Gstmt) Query(tid byte, where *gdo.WhereStmt) (string, error) {
	if table, ok := g.tables[tid]; ok {
		sql, err := table.BuildSelectStmt(where)
		if err != nil {
			return "", errors.Wrap(err, "生成 Select 語法時發生錯誤")
		}
		return sql, nil
	}
	return "", errors.New(fmt.Sprintf("Table %d is not exists.", tid))
}

// 取得符合 WhereStmt 條件的數據筆數
func (g *Gstmt) Count(tid byte, where *gdo.WhereStmt) (string, error) {
	if table, ok := g.tables[tid]; ok {
		sql := table.CountStmt(where)
		return sql, nil
	}
	return "", errors.New(fmt.Sprintf("Table(%d) is not exists.", tid))
}

// ====================================================================================================
// Update
// ====================================================================================================

func (g *Gstmt) Update(tid byte, data any, where *gdo.WhereStmt) (string, error) {
	if table, ok := g.tables[tid]; ok {
		table.InitByStruct(data)
		sql, err := table.Update(data, where)
		if err != nil {
			return "", errors.Wrap(err, "生成 Update 語法時發生錯誤")
		}
		return sql, nil
	}
	return "", errors.New(fmt.Sprintf("找不到編號為 %d 的表格", tid))
}

// ====================================================================================================
// Delete
// ====================================================================================================

func (g *Gstmt) DeleteBy(tid byte, where *gdo.WhereStmt) (string, error) {
	if table, ok := g.tables[tid]; ok {
		sql, err := table.BuildDeleteStmt(where)
		if err != nil {
			return "", err
		}
		return sql, nil
	}
	return "", errors.New(fmt.Sprintf("找不到編號為 %d 的表格", tid))
}
