package stmt

import (
	"fmt"
)

type DbOp byte

const (
	DbQuery DbOp = iota
	DbDelete
	DbInsert
	DbUpdate
)

type IStatement interface {
	ToStmt() string
}

// 預設參數
const (
	// INNODB 存儲引擎，多列索引的長度限制： 每個列的長度不能大於767 bytes；所有組成索引列的長度和不能大於 3072 bytes
	ENGINE string = "InnoDB"
	// 排序規則
	COLLATE string = "utf8mb4_bin"
	// Primary key 演算法
	ALGO string = "BTREE"
)

// 根據表格名稱，檢查該名稱的表格是否存在
func IsTableExists(tableName string) string {
	return fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE TABLE_NAME = '%s';", tableName)
}
