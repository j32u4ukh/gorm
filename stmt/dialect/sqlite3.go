package dialect

import (
	"fmt"

	"github.com/j32u4ukh/gorm/stmt/datatype"
)

type sqlite3 struct{}

func init() {
	RegisterDialect(SQLITE3, &sqlite3{})
}

// 變數類型，轉為 SQL 中的變數類型
func (s *sqlite3) TypeOf(dataType string) string {
	switch dataType {
	case datatype.TINYINT:
		fallthrough
	case datatype.SMALLINT:
		fallthrough
	case datatype.MEDIUMINT:
		fallthrough
	case datatype.INT:
		fallthrough
	case datatype.BIGINT:
		fallthrough
	case datatype.FLOAT:
		fallthrough
	case datatype.DOUBLE:
		fallthrough
	case datatype.DEMICAL:
		fallthrough
	case datatype.VARCHAR:
		fallthrough
	case datatype.CHAR:
		fallthrough
	case datatype.TINYTEXT:
		fallthrough
	case datatype.TEXT:
		fallthrough
	case datatype.MEDIUMTEXT:
		fallthrough
	case datatype.LONGTEXT:
		fallthrough
	case datatype.DATE:
		fallthrough
	case datatype.TIME:
		fallthrough
	case datatype.YEAR:
		fallthrough
	case datatype.DATETIME:
		fallthrough
	case datatype.TIMESTAMP:
		return dataType
	default:
		panic(fmt.Sprintf("Invalid variable type: %s.", dataType))
	}
}

// 根據 dataType 、當前的 size 以及 DB 本身的限制，對數值大小再定義
func (s *sqlite3) SizeOf(dataType string, size int32) int32 {
	switch dataType {
	case datatype.TINYINT:
		fallthrough
	case datatype.SMALLINT:
		fallthrough
	case datatype.MEDIUMINT:
		fallthrough
	case datatype.INT:
		fallthrough
	case datatype.BIGINT:
		fallthrough
	case datatype.FLOAT:
		fallthrough
	case datatype.DOUBLE:
		fallthrough
	case datatype.DEMICAL:
		fallthrough
	case datatype.VARCHAR:
		fallthrough
	case datatype.CHAR:
		fallthrough
	case datatype.TINYTEXT:
		fallthrough
	case datatype.TEXT:
		fallthrough
	case datatype.MEDIUMTEXT:
		fallthrough
	case datatype.LONGTEXT:
		fallthrough
	case datatype.DATE:
		fallthrough
	case datatype.TIME:
		fallthrough
	case datatype.YEAR:
		fallthrough
	case datatype.DATETIME:
		fallthrough
	case datatype.TIMESTAMP:
		return 0
	default:
		panic(fmt.Sprintf("Invalid variable type: %s.", dataType))
	}
}

// Protobuf 中的變數類型，轉為 SQL 中的變數類型
func (s *sqlite3) ProtoTypeOf(kind string) string {
	panic(fmt.Sprintf("invalid sql type %s", kind))
}

func (s *sqlite3) DbToProto(kind string) string {
	panic(fmt.Sprintf("invalid sql type %s", kind))
}

// 表格是否存在的 SQL 語法
func (s *sqlite3) IsTableExistsStmt(tableName string) string {
	return fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' and name = '%s';", tableName)
}

// 是否為數值類型
func (s *sqlite3) IsSortable(kind string) bool {
	return false
}

// 判斷變數類型(integer, float, text, ...)
func (s *sqlite3) GetKind(kind string) string {
	return kind
}
