package test

import (
	"fmt"
	"testing"

	"github.com/j32u4ukh/gorm/stmt"
	"github.com/j32u4ukh/gorm/stmt/datatype"
	"github.com/j32u4ukh/gorm/stmt/dialect"
)

func TestTableCreateStmt(t *testing.T) {
	answer := "CREATE TABLE IF NOT EXISTS `Desk` (`Id` INT(11) NOT NULL DEFAULT 0, `Content` VARCHAR(3000) NOT NULL DEFAULT '', PRIMARY KEY (`Id`) USING BTREE) ENGINE = InnoDB COLLATE = 'utf8mb4_bin';"
	tableParam := stmt.NewTableParam()
	tableParam.AddPrimaryKey("Id", "default")

	table := stmt.NewTable("Desk", tableParam, nil, stmt.ENGINE, stmt.COLLATE)
	col1 := stmt.NewColumnParam(1, "Id", datatype.INT, dialect.MARIA)
	col1.SetPrimaryKey("default")
	table.AddColumn(stmt.NewColumn(col1))

	col2 := stmt.NewColumnParam(2, "Content", datatype.VARCHAR, dialect.MARIA)
	table.AddColumn(stmt.NewColumn(col2))

	sql, err := table.BuildCreateStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestCreateStmt | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestCreateStmt |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestTableInsertStmt(t *testing.T) {
	answer := "INSERT INTO `Desk` (`Id`, `Content`) VALUES (41, ''), (42, 'not nil');"
	tableParam := stmt.NewTableParam()
	tableParam.AddPrimaryKey("Id", "default")
	var err error

	/////////////////////////////////////////////////////////////////////
	table := stmt.NewTable("Desk", tableParam, nil, stmt.ENGINE, stmt.COLLATE)
	col1 := stmt.NewColumnParam(1, "Id", datatype.INT, dialect.MARIA)
	col1.SetPrimaryKey("default")
	table.AddColumn(stmt.NewColumn(col1))

	col2 := stmt.NewColumnParam(2, "Content", datatype.VARCHAR, dialect.MARIA)
	// col2.SetCanNull(true)
	table.AddColumn(stmt.NewColumn(col2))
	/////////////////////////////////////////////////////////////////////
	table.Insert([]string{"41, ''"})
	table.Insert([]string{"42, 'not nil'"})

	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return
	}
	/////////////////////////////////////////////////////////////////////

	var sql string
	sql, err = table.BuildInsertStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestInsertStmt | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestInsertStmt |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestTableSelectStmt1(t *testing.T) {
	answer := "SELECT * FROM `Desk` ORDER BY `Id` DESC LIMIT 5 OFFSET 2;"
	tableParam := stmt.NewTableParam()
	tableParam.AddPrimaryKey("Id", "default")
	var err error

	/////////////////////////////////////////////////////////////////////
	table := stmt.NewTable("Desk", tableParam, nil, stmt.ENGINE, stmt.COLLATE)
	col1 := stmt.NewColumnParam(1, "Id", datatype.INT, dialect.MARIA)
	col1.SetPrimaryKey("default")
	table.AddColumn(stmt.NewColumn(col1))

	col2 := stmt.NewColumnParam(2, "Content", datatype.VARCHAR, dialect.MARIA)
	// col2.SetCanNull(true)
	table.AddColumn(stmt.NewColumn(col2))
	/////////////////////////////////////////////////////////////////////

	sql, err := table.
		SetOrderBy("Id").
		WhetherReverseOrder(true).
		SetLimit(5).
		SetOffset(2).
		ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestSelectStmt1 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestSelectStmt1 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestTableUpdateStmt(t *testing.T) {
	answer := "UPDATE `Desk` SET `Id` = 39, `Content` = 'Hello' WHERE `Id` = 39;"

	tableParam := stmt.NewTableParam()
	tableParam.AddPrimaryKey("Id", "default")
	var err error

	/////////////////////////////////////////////////////////////////////
	table := stmt.NewTable("Desk", tableParam, nil, stmt.ENGINE, stmt.COLLATE)
	col1 := stmt.NewColumnParam(1, "Id", datatype.INT, dialect.MARIA)
	col1.SetPrimaryKey("default")
	table.AddColumn(stmt.NewColumn(col1))

	col2 := stmt.NewColumnParam(2, "Content", datatype.VARCHAR, dialect.MARIA)
	// col2.SetCanNull(true)
	table.AddColumn(stmt.NewColumn(col2))
	/////////////////////////////////////////////////////////////////////

	table.SetUpdateCondition(stmt.WS().Eq("Id", "39"))

	sql, err := table.
		Update("Id", "39").
		Update("Content", "'Hello'").
		ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestUpdateStmt | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestUpdateStmt |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}
