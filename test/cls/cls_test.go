package main

import (
	"fmt"
	"testing"

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

func TestCreate(t *testing.T) {
	answer := "CREATE TABLE IF NOT EXISTS `demo`.`Desk` (`Id` INT(11) NOT NULL DEFAULT 0, `Content` VARCHAR(3000) NOT NULL DEFAULT '' COLLATE 'utf8mb4_bin', PRIMARY KEY (`Id`) USING BTREE) ENGINE = InnoDB COLLATE = 'utf8mb4_bin';"
	table := InitTable()
	sql, err := table.BuildCreateStmt()
	fmt.Printf("sql: %s\n", sql)

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestTableCreate | Errr: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestTableCreate |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestInsert(t *testing.T) {
	answer := "INSERT INTO `demo`.`Desk` (`Id`, `Content`) VALUES (3, 'content');"
	table := InitTable()
	desk := &Desk{Id: 3, Content: "content"}
	table.Insert([]any{desk})
	table.UseAntiInjection(true)
	sql, err := table.BuildInsertStmt()

	if err != nil || sql != answer {
		if err != nil {
			fmt.Printf("TestInsert | Error: %+v\n", err)
		}

		if sql != answer {
			fmt.Printf("TestInsert |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestQuery(t *testing.T) {
	answer := "SELECT * FROM `demo`.`Desk` WHERE `Id` = 3;"
	table := InitTable()
	// desk := &Desk{Id: 3, Content: "content"}
	table.UseAntiInjection(true)
	sql, err := table.BuildSelectStmt(gdo.WS().Eq("Id", 3))

	if err != nil || sql != answer {
		if err != nil {
			fmt.Printf("TestInsert | Error: %+v\n", err)
		}

		if sql != answer {
			fmt.Printf("TestInsert |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestUpdate(t *testing.T) {
	answer := "UPDATE `demo`.`Desk` SET `Id` = 3, `Content` = 'content' WHERE `Id` = 3;"
	table := InitTable()
	desk := &Desk{Id: 3, Content: "content"}
	table.UseAntiInjection(true)
	sql, err := table.Update(desk, gdo.WS().Eq("Id", 3))

	if err != nil || sql != answer {
		if err != nil {
			fmt.Printf("TestInsert | Error: %+v\n", err)
		}

		if sql != answer {
			fmt.Printf("TestInsert |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestTableDelete(t *testing.T) {
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
