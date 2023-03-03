package dialect

import (
	"fmt"
	"strings"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gorm/stmt/datatype"
)

type mysql struct {
	KindMap map[string]*cntr.Array[string]
}

func init() {
	RegisterDialect(MYSQL, &mysql{
		KindMap: map[string]*cntr.Array[string]{
			"INTEGER": cntr.NewArray(
				"TINYINT",
				"SMALLINT",
				"MEDIUMINT",
				"INT",
				"BIGINT",
			),
			"FLOAT": cntr.NewArray(
				"FLOAT",
				"DOUBLE",
				"DEMICAL",
			),
			"STRING": cntr.NewArray(
				"VARCHAR",
				"CHAR",
				"TINYTEXT",
				"TEXT",
				"MEDIUMTEXT",
				"LONGTEXT",
			),
		},
	})
}

// 變數類型，轉為 SQL 中的變數類型
func (s *mysql) TypeOf(dataType string) string {
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
func (s *mysql) SizeOf(dataType string, size int32) int32 {
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
func (s *mysql) ProtoTypeOf(kind string) string {
	switch kind {
	case "INT32":
		return "INT"
	case "INT64":
		return "BIGINT"
	case "BOOL":
		return "TINYINT"
	case "STRING":
		return "VARCHAR"
	// 原生 SQL 變數類型，無須修改
	default:
		return kind
	}
}

func (s *mysql) DbToProto(kind string) string {
	switch kind {
	case "INT":
		return "int32"
	case "BIGINT":
		return "int64"
	case "TINYINT":
		return "bool"
	default:
		return "string"
	}
}

// 表格是否存在的 SQL 語法
func (s *mysql) IsTableExistsStmt(tableName string) string {
	return fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE TABLE_NAME = '%s';", tableName)
}

// 是否可適用排序演算法
func (s *mysql) IsSortable(kind string) bool {
	return s.KindMap["STRING"].Contains(strings.ToUpper(kind))
}

// 判斷變數類型(integer, float, text, ...)
func (s *mysql) GetKind(kind string) string {
	for k, v := range s.KindMap {
		if v.Contains(kind) {
			return k
		}
	}
	return "Unknown"
}
