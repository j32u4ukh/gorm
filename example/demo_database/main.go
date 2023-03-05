package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gorm/database"
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
	case "d":
		DeleteDemo()
	default:
		fmt.Printf("No invalid command(%s).\n", command)
	}
}

func CreateDemo() {
	sql := "CREATE TABLE `Desk` ( `index` INT(11) NOT NULL AUTO_INCREMENT COMMENT '索引值', `user_name` VARCHAR(20) NOT NULL COMMENT '玩家名稱' COLLATE 'utf8mb4_bin',`item_id` INT(11) NOT NULL COMMENT '物品 ID',`time` TIMESTAMP NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Log 建立時間',PRIMARY KEY (`index`) USING BTREE) COLLATE='utf8mb4_bin' ENGINE=InnoDB;"
	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}

func InsertDemo() {
	sql := "INSERT INTO `pekomiko`.`Desk` (`user_name`, `item_id`) VALUES ('sss', '23');"

	for i := 0; i < 10; i++ {
		result, err := db.Exec(sql)

		if err != nil {
			fmt.Printf("Insert err: %+v\n", err)
			return
		}

		fmt.Printf("result: %s\n", result)
	}
}

func QueryDemo() {
	sql := "SELECT `index`, `user_name`, `item_id`, `time` FROM `pekomiko`.`Desk` LIMIT 10;"
	result, err := db.Query(sql)

	if err != nil {
		fmt.Printf("Select err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
	for i, data := range result.Datas {
		fmt.Printf("i: %d, data: %s\n", i, cntr.SliceToString(data))
	}
}

func UpdateDemo() {
	sql := "UPDATE `pekomiko`.`Desk` SET `item_id`=24 WHERE `index`=1;"
	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("Update err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}

func DeleteDemo() {
	sql := "DELETE FROM `pekomiko`.`Desk` WHERE `item_id`=24;"
	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("Update err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}
