package main

import (
	"fmt"

	"github.com/j32u4ukh/gorm/cls"
	"github.com/j32u4ukh/gorm/gdo"
	"github.com/j32u4ukh/gorm/stmt/dialect"
)

type Desk struct {
	Id      int    `gorm:"pk=default"`
	Content string `gorm:"size=3000"`
}

func InitTable() *cls.StructTable {
	desk := &Desk{}
	tableParams, columnParams, err := cls.GetParams(desk, dialect.MARIA)
	table := cls.NewStructTable("Desk", tableParams, columnParams, dialect.MARIA)
	table.SetDbName("demo")
	table.InitByStruct(desk)
	if err != nil {
		fmt.Printf("BuildCreateStmt err: %+v\n", err)
		return nil
	}
	return table
}

func main() {
	answer := "DELETE FROM `demo`.`Desk` WHERE `Id` = 3;"
	table := InitTable()
	// desk := &Desk{Id: 3, Content: "content"}
	table.UseAntiInjection(true)
	sql, err := table.BuildDeleteStmt(gdo.WS().Eq("Id", 3))

	if err != nil || sql != answer {
		if err != nil {
			fmt.Printf("TestInsert | Error: %+v\n", err)
		}

		if sql != answer {
			fmt.Printf("TestInsert |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}
