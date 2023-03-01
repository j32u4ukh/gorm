package test

import (
	"fmt"
	"testing"

	"github.com/j32u4ukh/gorm/stmt"
	"github.com/j32u4ukh/gorm/stmt/datatype"
	"github.com/j32u4ukh/gorm/stmt/dialect"
)

func TestCreateStmt(t *testing.T) {
	answer := "CREATE TABLE `Desk` (`Id` INT(11) NOT NULL, `Content` VARCHAR(3000) NOT NULL, PRIMARY KEY (`Id`) USING BTREE) ENGINE = InnoDB COLLATE = 'utf8mb4_bin';"
	tableParam := stmt.NewTableParam()
	tableParam.AddPrimaryKey("Id", "default")

	cs := stmt.NewCreateStmt("Desk", tableParam, nil, stmt.ENGINE, stmt.COLLATE)

	col1 := stmt.NewColumnParam(1, "Id", datatype.INT, dialect.MARIA)
	col1.SetPrimaryKey("default")
	cs.AddColumn(stmt.NewColumn(col1))

	col2 := stmt.NewColumnParam(2, "Content", datatype.VARCHAR, dialect.MARIA)
	// col2.SetCanNull(true)
	cs.AddColumn(stmt.NewColumn(col2))

	sql, err := cs.ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestCreateStmt | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestCreateStmt |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestTableCreateStmt(t *testing.T) {
	answer := "CREATE TABLE `Desk` (`Id` INT(11) NOT NULL, `Content` VARCHAR(3000) NOT NULL, PRIMARY KEY (`Id`) USING BTREE) ENGINE = InnoDB COLLATE = 'utf8mb4_bin';"
	tableParam := stmt.NewTableParam()
	tableParam.AddPrimaryKey("Id", "default")

	table := stmt.NewTable("Desk", tableParam, nil, stmt.ENGINE, stmt.COLLATE)
	col1 := stmt.NewColumnParam(1, "Id", datatype.INT, dialect.MARIA)
	col1.SetPrimaryKey("default")
	table.AddColumn(stmt.NewColumn(col1))

	col2 := stmt.NewColumnParam(2, "Content", datatype.VARCHAR, dialect.MARIA)
	// col2.SetCanNull(true)
	table.AddColumn(stmt.NewColumn(col2))

	sql, err := table.BuildCreateStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestCreateStmt | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestCreateStmt |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestInsertStmt(t *testing.T) {
	answer := "INSERT INTO `Desk` (`Id`, `Content`) VALUES (41, NULL), (42, 'not nil');"
	is := stmt.NewInsertStmt("Desk")
	is.SetColumnNames([]string{"Id", "Content"})
	sql, err := is.Insert([]string{"41", "NULL"}).
		Insert([]string{"42", "'not nil'"}).
		ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestInsertStmt | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestInsertStmt |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestTableInsertStmt(t *testing.T) {
	answer := "INSERT INTO `Desk` (`Id`, `Content`) VALUES (41, NIL), (42, 'not nil');"
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
	err = table.InsertWithCheck(map[string]string{"Id": "41"})

	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return
	}

	err = table.InsertWithCheck(map[string]string{"Id": "42", "Content": "'not nil'"})

	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return
	}
	/////////////////////////////////////////////////////////////////////

	var sql string
	sql, err = table.BuildInsertStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestInsertStmt | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestInsertStmt |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestSelectStmt1(t *testing.T) {
	answer := "SELECT * FROM `Desk` ORDER BY `Id` DESC LIMIT 5 OFFSET 2;"
	sql, err := stmt.NewSelectStmt("Desk").
		SetOrderBy("Id").
		WhetherReverseOrder(true).
		SetLimit(5).
		SetOffset(2).
		ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestSelectStmt1 | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestSelectStmt1 |\nanswer: %s\nsql: %s", answer, sql))
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
			t.Error(fmt.Sprintf("TestSelectStmt1 | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestSelectStmt1 |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestSelectStmt2(t *testing.T) {
	answer := "SELECT * FROM `customers` WHERE `Name` = '王二';"
	where := stmt.WS().Eq("Name", "'王二'")
	sql, err := stmt.NewSelectStmt("customers").
		SetCondition(where).
		ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestSelectStmt2 | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestSelectStmt2 |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestSelectStmt3(t *testing.T) {
	answer := "SELECT Name, Phone FROM `customers` WHERE `City` = '台北市' AND `Salary` >= 50000;"
	where := stmt.WS().
		AddAndCondtion(stmt.WS().Eq("City", "'台北市'")).
		AddAndCondtion(stmt.WS().Ge("Salary", "50000"))

	sql, err := stmt.
		NewSelectStmt("customers").
		Query("Name").
		Query("Phone").
		SetCondition(where).
		ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestSelectStmt3 | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestSelectStmt3 |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestSelectStmt4(t *testing.T) {
	answer := "SELECT Name, Phone FROM `customers` WHERE `Name` = 'Sam' OR (`City` = '台北市' AND `Salary` >= 50000);"

	where := stmt.WS().
		AddOrCondtion(stmt.WS().Eq("Name", "'Sam'")).
		AddOrCondtion(stmt.WS().
			AddAndCondtion(stmt.WS().Eq("City", "'台北市'")).
			AddAndCondtion(stmt.WS().Ge("Salary", "50000")).
			SetBrackets())

	sql, err := stmt.
		NewSelectStmt("customers").
		Query("Name").
		Query("Phone").
		SetCondition(where).
		ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestSelectStmt4 | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestSelectStmt4 |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestSelectStmt5(t *testing.T) {
	answer := "SELECT * FROM `emp` WHERE NOT (`Salary` > 50000);"

	where := stmt.WS().Gt("Salary", "50000").SetNotCondition()

	sql, err := stmt.
		NewSelectStmt("emp").
		SetCondition(where).
		ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestSelectStmt5 | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestSelectStmt5 |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestSelectStmt6(t *testing.T) {
	answer := "SELECT * FROM `emp` WHERE (`comm` IS NULL) AND (`id` IS NOT NULL);"

	where := stmt.WS().
		AddAndCondtion(stmt.WS().CheckNull("comm", true).SetBrackets()).
		AddAndCondtion(stmt.WS().CheckNull("id", false).SetBrackets())

	sql, err := stmt.
		NewSelectStmt("emp").
		SetCondition(where).
		ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestSelectStmt6 | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestSelectStmt6 |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestSelectStmt7(t *testing.T) {
	answer := "SELECT * FROM `BI.common_log_2022_09_22_00_00` WHERE `bi_type` IN (10, 11, 12) AND `create_time` BETWEEN '2022-09-22 00:00:00' AND '2022-09-22 23:59:59' AND `d_1` != '';"

	where := stmt.WS().
		AddAndCondtion(stmt.WS().In("bi_type", "10", "11", "12")).
		AddAndCondtion(stmt.WS().Between("create_time", "'2022-09-22 00:00:00'", "'2022-09-22 23:59:59'")).
		AddAndCondtion(stmt.WS().Ne("d_1", "''"))

	sql, err := stmt.
		NewSelectStmt("BI.common_log_2022_09_22_00_00").
		SetCondition(where).
		ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestSelectStmt7 | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestSelectStmt7 |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestUpdateStmt(t *testing.T) {
	answer := "UPDATE `Desk` SET `Id` = 39, `Content` = 'Hello' WHERE `Id` = 39;"

	update := stmt.NewUpdateStmt("Desk")
	update.SetCondition(stmt.WS().Eq("Id", "39"))
	sql, err := update.
		Update("Id", "39").
		Update("Content", "'Hello'").
		ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestUpdateStmt | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestUpdateStmt |\nanswer: %s\nsql: %s", answer, sql))
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
			t.Error(fmt.Sprintf("TestUpdateStmt | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestUpdateStmt |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestDeleteStmt(t *testing.T) {
	answer := "DELETE FROM `demo2`.`Desk` WHERE `Id` = 39;"

	tableParam := stmt.NewTableParam()
	tableParam.AddPrimaryKey("Id", "default")
	var err error

	/////////////////////////////////////////////////////////////////////
	table := stmt.NewTable("Desk", tableParam, nil, stmt.ENGINE, stmt.COLLATE)
	table.SetDbName("demo2")

	col1 := stmt.NewColumnParam(1, "Id", datatype.INT, dialect.MARIA)
	col1.SetPrimaryKey("default")
	table.AddColumn(stmt.NewColumn(col1))

	col2 := stmt.NewColumnParam(2, "Content", datatype.VARCHAR, dialect.MARIA)
	// col2.SetCanNull(true)
	table.AddColumn(stmt.NewColumn(col2))
	/////////////////////////////////////////////////////////////////////

	table.SetDeleteCondition(stmt.WS().Eq("Id", "39"))
	sql, err := table.BuildDeleteStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestDeleteStmt | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestDeleteStmt |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

// func TestBatchUpdateStmt(t *testing.T) {
// 	answer := "UPDATE Desk SET `Value` = CASE `Id` WHEN 1 THEN '3' WHEN 3 THEN '11' END, `Weight` = CASE `Id` WHEN 1 THEN 39 WHEN 2 THEN 79 END WHERE `Id` IN (1, 2, 3);"
//
// 	sql, err := stmt.
// 		NewBatchUpdateStmt("Desk", "Id").
// 		AddData(map[string]any{"Id": 1, "Value": "3", "Weight": 39}).
// 		AddData(map[string]any{"Id": 2, "Weight": 79}).
// 		AddData(map[string]any{"Id": 3, "Value": "11"}).
// 		ToStmt()

// 	if err != nil || sql != answer {
// 		if err != nil {
// 			t.Error(fmt.Sprintf("TestBatchUpdateStmt | Error: %+v\n", err))
// 		}

// 		if sql != answer {
// 			t.Error(fmt.Sprintf("TestBatchUpdateStmt |\nanswer: %s\nsql: %s", answer, sql))
// 		}
// 	}
// }
