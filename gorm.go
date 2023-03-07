package gorm

import (
	"fmt"

	"github.com/j32u4ukh/gorm/cls"
	"github.com/j32u4ukh/gorm/gdo"
	"github.com/j32u4ukh/gorm/stmt"
	"github.com/j32u4ukh/gorm/stmt/dialect"
	"github.com/pkg/errors"
)

const (
	TIME_LAYOUT string = "2006-01-02 15:04:05"
)

type ITable interface {
	SetDbName(dbName string)
	UseAntiInjection(active bool)
	InitByStruct(data any)
	// Create
	BuildCreateStmt() (string, error)
	// Insert
	Insert(datas []any) error
	BuildInsertStmt() (string, error)
	// Select
	Query(column string) *stmt.SelectStmt
	CountMode() *stmt.SelectStmt
	SetOrderBy(column string) *stmt.SelectStmt
	SetSelectCondition(where *gdo.WhereStmt)
	WhetherReverseOrder(reverse bool) *stmt.SelectStmt
	SetLimit(limit int32) *stmt.SelectStmt
	SetOffset(offset int32) *stmt.SelectStmt
	BuildSelectStmt() (string, error)
	// Update
	Update(data any, where *gdo.WhereStmt)
	BuildUpdateStmt() (string, error)
	// Delete
	SetDeleteCondition(where *gdo.WhereStmt)
	BuildDeleteStmt(where *gdo.WhereStmt) (string, error)
}

var gstmtMap map[byte]*Gorm

type Gorm struct {
	DbName string
	Dial   dialect.SQLDialect
	tables map[byte]ITable
}

func init() {
	gstmtMap = map[byte]*Gorm{}
}

func SetGorm(index byte, dbName string, dialect dialect.SQLDialect) (*Gorm, error) {
	if _, ok := gstmtMap[index]; ok {
		return nil, errors.New(fmt.Sprintf("index: %d 已使用", index))
	}
	g, err := NewGorm(dbName, dialect)
	if err != nil {
		return nil, errors.Wrap(err, "建立 Gorm 時發生錯誤")
	}
	gstmtMap[index] = g
	return g, nil
}

func GetGorm(index byte) (*Gorm, error) {
	if g, ok := gstmtMap[index]; ok {
		return g, nil
	}
	return nil, errors.New(fmt.Sprintf("index: %d 沒有對應的 *Gorm", index))
}

// 由外部定義實際的 DB 實體，GoSql 則對該 DB 實體進行操作
// db: DB 指標
// dialect: 資料庫語言(mysql/maria/...)
func NewGorm(dbName string, dial dialect.SQLDialect) (*Gorm, error) {
	g := &Gorm{
		DbName: dbName,
		Dial:   dial,
		tables: map[byte]ITable{},
	}
	return g, nil
}

func (g *Gorm) SetDbName(dbName string) {
	g.DbName = dbName
	for _, t := range g.tables {
		t.SetDbName(dbName)
	}
}

func (g *Gorm) UseAntiInjection(active bool) {
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
func (g *Gorm) SetTable(tid byte, tableName string, dial dialect.SQLDialect, data any) error {
	// 檢查表格是否已存在變數中
	if _, ok := g.tables[tid]; ok {
		return errors.New(fmt.Sprintf("表格 %s 已存在", tableName))
	}

	// 根據檔案名稱，取得表格與欄位參數
	tableParam, colParams, err := cls.GetParams(data, dial)

	if err != nil {
		return errors.Wrap(err, "讀取參數時發生錯誤")
	}

	t := cls.NewStructTable(tableName, tableParam, colParams, g.Dial)
	t.InitByStruct(data)
	t.SetDbName(g.DbName)

	// 將 Table 加入管理，可利用 tableName 進行存取
	g.tables[tid] = t
	fmt.Printf("(g *Gorm) SetTable | table: %+v.\n", t)
	return nil
}

func (g *Gorm) GetTable(tid byte) ITable {
	if table, ok := g.tables[tid]; ok {
		return table
	}
	return nil
}

// // func (g *Gorm) DeleteTable(tid byte)

// // ====================================================================================================
// // Insert
// // ====================================================================================================
// func (g *Gorm) Insert(tid byte, datas []any) (string, error) {
// 	if table, ok := g.tables[tid]; ok {
// 		table.InitByStruct(datas[0])
// 		table.Insert(datas)
// 		sql, err := table.BuildInsertStmt()
// 		if err != nil {
// 			return "", errors.Wrap(err, "Failed to build InsertStmt.")
// 		}
// 		return sql, nil
// 	}
// 	return "", errors.New(fmt.Sprintf("找不到編號為 %d 的表格", tid))
// }

// // ====================================================================================================
// // Select
// // ====================================================================================================

// func (g *Gorm) Query(tid byte, where *gdo.WhereStmt) (string, error) {
// 	if table, ok := g.tables[tid]; ok {
// 		table.SetSelectCondition(where)
// 		sql, err := table.BuildSelectStmt()
// 		if err != nil {
// 			return "", errors.Wrap(err, "生成 Select 語法時發生錯誤")
// 		}
// 		return sql, nil
// 	}
// 	return "", errors.New(fmt.Sprintf("Table %d is not exists.", tid))
// }

// // 取得符合 WhereStmt 條件的數據筆數
// func (g *Gorm) Count(tid byte, where *gdo.WhereStmt) (string, error) {
// 	if table, ok := g.tables[tid]; ok {
// 		sql := table.CountStmt(where)
// 		return sql, nil
// 	}
// 	return "", errors.New(fmt.Sprintf("Table(%d) is not exists.", tid))
// }

// // ====================================================================================================
// // Update
// // ====================================================================================================

// func (g *Gorm) Update(tid byte, data any, where *gdo.WhereStmt) (string, error) {
// 	if table, ok := g.tables[tid]; ok {
// 		table.InitByStruct(data)
// 		table.Update(data, where)
// 		sql, err := table.BuildUpdateStmt()
// 		if err != nil {
// 			return "", errors.Wrap(err, "生成 Update 語法時發生錯誤")
// 		}
// 		return sql, nil
// 	}
// 	return "", errors.New(fmt.Sprintf("找不到編號為 %d 的表格", tid))
// }

// // ====================================================================================================
// // Delete
// // ====================================================================================================

// func (g *Gorm) DeleteBy(tid byte, where *gdo.WhereStmt) (string, error) {
// 	if table, ok := g.tables[tid]; ok {
// 		sql, err := table.BuildDeleteStmt(where)
// 		if err != nil {
// 			return "", err
// 		}
// 		return sql, nil
// 	}
// 	return "", errors.New(fmt.Sprintf("找不到編號為 %d 的表格", tid))
// }
