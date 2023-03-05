package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gorm/database"
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

func CreateDemo() {
	tableParam := stmt.NewTableParam()
	tableParam.AddPrimaryKey("Id", "default")

	cs := stmt.NewCreateStmt("Desk", tableParam, nil, stmt.ENGINE, stmt.COLLATE)

	col1 := stmt.NewColumnParam(1, "Id", datatype.INT, dialect.MARIA)
	col1.SetPrimaryKey("default")
	cs.AddColumn(stmt.NewColumn(col1))

	col2 := stmt.NewColumnParam(2, "Content", datatype.VARCHAR, dialect.MARIA)
	cs.AddColumn(stmt.NewColumn(col2))

	sql, err := cs.ToStmt()

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

func InsertDemo() {
	is := stmt.NewInsertStmt("Desk")
	is.SetColumnNames([]string{"Id", "Content"})
	for i := 1; i <= 10; i++ {
		is.Insert([]string{fmt.Sprintf("%d", 40+i), fmt.Sprintf("'content%d'", i)})
	}
	sql, err := is.ToStmt()

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
	sql, err := stmt.NewSelectStmt("Desk").
		SetOrderBy("Id").
		WhetherReverseOrder(true).
		SetLimit(5).
		SetOffset(2).
		ToStmt()

	if err != nil {
		return
	}

	fmt.Printf("QueryDemo | sql: %s\n", sql)
	result, err := db.Query(sql)

	if err != nil {
		fmt.Printf("Query err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
	for i, data := range result.Datas {
		fmt.Printf("i: %d, data: %s\n", i, cntr.SliceToString(data))
	}
}

func UpdateDemo() {
	update := stmt.NewUpdateStmt("Desk")
	update.SetCondition(stmt.WS().Eq("Id", "42"))
	sql, err := update.
		Update("Id", "42").
		Update("Content", "'Hello'").
		ToStmt()

	if err != nil {
		return
	}

	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("Update err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}

func BatchUpdateDemo() {
	bus := stmt.NewBatchUpdateStmt("Desk", "Id")
	bus.Update("41", "Content", "3")
	bus.Update("43", "Content", "4")
	bus.Update("45", "Content", "6")
	sql, err := bus.ToStmt()

	if err != nil {
		return
	}

	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("BatchUpdate err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}

func DeleteDemo() {
	del := stmt.NewDeleteStmt("Desk")
	del.SetDbName("pekomiko")
	del.SetCondition(stmt.WS().Eq("Id", "49"))
	sql, err := del.ToStmt()

	if err != nil {
		return
	}

	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("Delete err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}
