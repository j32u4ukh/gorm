package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gorm"
	"github.com/j32u4ukh/gorm/database"
	"github.com/j32u4ukh/gorm/gdo"
	"github.com/j32u4ukh/gorm/stmt/dialect"
)

const GIndex byte = 0
const TID byte = 0

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

	gorm.SetGorm(GIndex, dc.Name, dialect.MARIA)

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

type Desk struct {
	Id      int    `gorm:"pk=default"`
	Content string `gorm:"size=3000"`
}

func InitTable() *gorm.Gorm {
	desk := &Desk{}
	g, _ := gorm.GetGorm(GIndex)
	err := g.SetTable(TID, "Desk", dialect.MARIA, desk)
	if err != nil {
		fmt.Printf("SetTable err: %+v\n", err)
		return nil
	}
	table := g.GetTable(TID)
	table.SetDbName("pekomiko")
	table.UseAntiInjection(true)
	table.InitByStruct(desk)
	return g
}

func CreateDemo() {
	g := InitTable()
	sql, err := g.GetTable(TID).BuildCreateStmt()

	if err != nil {
		return
	}

	// fmt.Printf("CreateDemo | sql: %s\n", sql)

	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("Create err: %+v\nsql: %s\n", err, sql)
		return
	}

	fmt.Printf("result: %s\n", result)
}

func InsertDemo() {
	g := InitTable()
	table := g.GetTable(TID)

	for i := 1; i <= 10; i++ {
		desk := &Desk{Id: 40 + i, Content: fmt.Sprintf("content%d", i)}
		table.Insert([]any{desk})
	}

	sql, err := table.BuildInsertStmt()

	if err != nil {
		return
	}

	// fmt.Printf("InsertDemo | sql: %s\n", sql)
	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("Insert err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}

func QueryDemo() {
	g := InitTable()
	table := g.GetTable(TID)
	table.SetOrderBy("Id").
		WhetherReverseOrder(true).
		SetLimit(5).
		SetOffset(2)
	sql, err := table.BuildSelectStmt()

	if err != nil {
		return
	}

	// fmt.Printf("QueryDemo | sql: %s\n", sql)
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
	g := InitTable()
	table := g.GetTable(TID)
	desk := &Desk{Id: 47, Content: "content"}
	table.Update(desk, gdo.WS().Eq("Id", 47))
	sql, err := table.BuildUpdateStmt()

	if err != nil {
		return
	}

	// fmt.Printf("UpdateDemo | sql: %s\n", sql)
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
	g := InitTable()
	table := g.GetTable(TID)
	table.SetDeleteCondition(gdo.WS().Eq("Id", 49))
	sql, err := table.BuildDeleteStmt(nil)

	if err != nil {
		return
	}

	// fmt.Printf("DeleteDemo | sql: %s\n", sql)
	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}
