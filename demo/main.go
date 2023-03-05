package main

import (
	"fmt"

	"github.com/j32u4ukh/gorm/gdo"
	"github.com/j32u4ukh/gorm/stmt"
	"github.com/j32u4ukh/gorm/stmt/datatype"
	"github.com/j32u4ukh/gorm/stmt/dialect"
)

type Desk struct {
	Id      int    `gorm:"pk=default"`
	Content string `gorm:"size=3000"`
}

func InitTable() *gdo.Table {
	tableName := "StmtDesk"
	tableParam := stmt.NewTableParam()

	// NewTable(name string, tableParam *stmt.TableParam, columnParams []*stmt.ColumnParam, engine string, collate string, dial string)
	table := gdo.NewTable(tableName, tableParam, nil, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)
	table.SetDbName("demo2")
	table.UseAntiInjection(true)
	col1 := stmt.NewColumnParam(1, "Id", datatype.INT, dialect.MARIA)
	col1.SetPrimaryKey("default")
	table.AddColumn(stmt.NewColumn(col1))

	col2 := stmt.NewColumnParam(2, "Content", datatype.VARCHAR, dialect.MARIA)
	// col2.SetCanNull(true)
	table.AddColumn(stmt.NewColumn(col2))
	return table
}

func main() {
	answer := "INSERT INTO `demo2`.`StmtDesk` (`Id`, `Content`) VALUES (50, '\\'); SELECT * FROM `demo2`.`StmtDesk`; -- hack');"
	table := InitTable()
	// map[string]any{"Id": 50, "Content": "'); SELECT * FROM `demo2`.`StmtDesk`; -- hack"}
	table.Insert([]any{50, "'); SELECT * FROM `demo2`.`StmtDesk`; -- hack"}, nil)
	sql, err := table.BuildInsertStmt()

	if err != nil {
		fmt.Printf("BuildCreateStmt err: %+v\n", err)
		return
	}

	if answer != sql {
		fmt.Printf("answer: %s\nsql***: %s\n", answer, sql)
	} else {
		fmt.Printf("sql: %s\n", sql)
	}
}
