package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gorm/database"
	"github.com/j32u4ukh/gorm/gdo"
	"github.com/j32u4ukh/gorm/stmt"
	"github.com/j32u4ukh/gorm/stmt/datatype"
	"github.com/j32u4ukh/gorm/stmt/dialect"
)

var db *database.Database

func main() {
	command := strings.ToLower(os.Args[1])
	conf, err := database.NewConfig("../config/config.yaml")

	if err != nil {
		fmt.Printf("讀取 Config 檔時發生錯誤, err: %+v\n", err)
		return
	}
	dc := conf.GetDatabase()
	db, err = database.Connect(0, dc.UserName, dc.Password, dc.Server, dc.Port, dc.Name)
	if err != nil {
		fmt.Printf("與資料庫連線時發生錯誤, err: %+v\n", err)
		return
	}
	defer db.Close()

	switch command {
	case "c":
		CreateDemo()
	case "i":
		InsertDemo()
	case "q":
		QueryDemo()
	case "u":
		UpdateDemo()
	case "bu":
		BatchUpdateDemo()
	case "d":
		DeleteDemo()
	default:
		fmt.Printf("No invalid command(%s).\n", command)
	}
}

func InitTable() *gdo.Table {
	tableName := "Desk"
	tableParam := stmt.NewTableParam()

	// NewTable(name string, tableParam *stmt.TableParam, columnParams []*stmt.ColumnParam, engine string, collate string, dial string)
	table := gdo.NewTable(tableName, tableParam, nil, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)
	table.SetDbName("pekomiko")
	table.UseAntiInjection(true)
	col1 := stmt.NewColumnParam(1, "Id", datatype.INT, dialect.MARIA)
	col1.SetPrimaryKey("default")
	table.AddColumn(stmt.NewColumn(col1))

	col2 := stmt.NewColumnParam(2, "Content", datatype.VARCHAR, dialect.MARIA)
	// col2.SetCanNull(true)
	table.AddColumn(stmt.NewColumn(col2))
	return table
}

func CreateDemo() {
	table := InitTable()
	sql, err := table.BuildCreateStmt()

	if err != nil {
		return
	}

	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("Create err: %+v\nsql: %s\n", err, sql)
		return
	}

	fmt.Printf("result: %s\n", result)
}

func InsertDemo() {
	table := InitTable()
	for i := 1; i <= 10; i++ {
		table.Insert([]any{40 + i, fmt.Sprintf("'content%d'", i)}, nil)
	}
	sql, err := table.BuildInsertStmt()

	if err != nil {
		return
	}

	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("Insert err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}

func QueryDemo() {
	table := InitTable()
	table.SetOrderBy("Id").
		WhetherReverseOrder(true).
		SetLimit(5).
		SetOffset(2)

	sql, err := table.BuildSelectStmt()

	if err != nil {
		return
	}

	fmt.Printf("QueryDemo | sql: %s\n", sql)
	result, err := db.Query(sql)

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
	for i, data := range result.Datas {
		fmt.Printf("i: %d, data: %s\n", i, cntr.SliceToString(data))
	}
}

func UpdateDemo() {
	table := InitTable()
	table.SetUpdateCondition(gdo.WS().Eq("Id", 42))
	table.Update("Content", "'Hello'", nil)
	sql, err := table.BuildUpdateStmt()

	if err != nil {
		return
	}

	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}

func BatchUpdateDemo() {
	// bus := stmt.NewBatchUpdateStmt("Desk", "Id")
	// bus.Update("41", "Content", "3")
	// bus.Update("43", "Content", "4")
	// bus.Update("45", "Content", "6")
	// sql, err := bus.ToStmt()

	// if err != nil {
	// 	return
	// }

	// result, err := db.Exec(sql)

	// if err != nil {
	// 	fmt.Printf("Create err: %+v\n", err)
	// 	return
	// }

	// fmt.Printf("result: %s\n", result)
}

func DeleteDemo() {
	table := InitTable()
	table.SetDeleteCondition(gdo.WS().Eq("Id", 49))
	sql, err := table.BuildDeleteStmt()

	if err != nil {
		return
	}

	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}
