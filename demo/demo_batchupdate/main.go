package main

import (
	"fmt"

	"github.com/j32u4ukh/gorm/stmt"
)

type Desk struct {
	Id     int
	Length int
	Width  int
}

func main() {
	bus := stmt.NewBatchUpdateStmt("Desk", "id")
	bus.Update("0", "length", "3")
	bus.Update("0", "width", "5")
	bus.Update("1", "length", "4")
	bus.Update("1", "width", "6")
	bus.Update("2", "length", "6")
	bus.Update("2", "width", "5")
	sql, err := bus.ToStmt()
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return
	}
	/* UPDATE Desk
	SET
	`length` = CASE `id`
	WHEN 0 THEN 3
	WHEN 1 THEN 4
	WHEN 2 THEN 6
	END,
	`width` = CASE `id`
	WHEN 0 THEN 5
	WHEN 1 THEN 6
	WHEN 2 THEN 5
	END
	WHERE `id` IN (0, 1, 2);*/
	fmt.Printf("sql: %s\n", sql)
}
