package dialect

import (
	"fmt"
	"strings"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gorm/stmt/datatype"
)

type maria struct {
	KindMap map[string]*cntr.Array[string]
}

func init() {
	RegisterDialect(MARIA, &maria{
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
			"TIME": cntr.NewArray(
				"DATE",
				"TIME",
				"YEAR",
				"DATETIME",
				"TIMESTAMP",
			),
		},
	})
}

// 變數類型，轉為 SQL 中的變數類型
func (s *maria) TypeOf(dataType string) string {
	switch dataType {
	case datatype.TINYINT, datatype.SMALLINT, datatype.MEDIUMINT, datatype.INT, datatype.BIGINT,
		datatype.FLOAT, datatype.DOUBLE, datatype.DEMICAL,
		datatype.VARCHAR, datatype.CHAR, datatype.TINYTEXT, datatype.TEXT, datatype.MEDIUMTEXT, datatype.LONGTEXT,
		datatype.DATE, datatype.TIME, datatype.YEAR, datatype.DATETIME, datatype.TIMESTAMP:
		return dataType
	case datatype.BOOL:
		return datatype.TINYINT
	default:
		// panic(fmt.Sprintf("Invalid variable type: %s.", dataType))
		fmt.Printf("(s *maria) TypeOf | dataType: %s\n", dataType)
		return datatype.VARCHAR
	}
}

// 根據 dataType 、當前的 size 以及 DB 本身的限制，對數值大小再定義
func (s *maria) SizeOf(dataType string, size int32) int32 {
	if size <= 0 {
		switch dataType {
		case datatype.BOOL:
			return 1
		case datatype.TINYINT:
			return 4
		case datatype.SMALLINT:
			return 6
		case datatype.MEDIUMINT:
			return 9
		case datatype.INT:
			return 11
		case datatype.BIGINT:
			return 20
		case datatype.VARCHAR:
			return 3000
		case datatype.CHAR:
			return 50
		}
	}

	switch dataType {
	// ==================================================
	// DB 本身有其預設值，無法自定義大小的類型，一律回傳 0
	case datatype.TIMESTAMP, datatype.DOUBLE, datatype.FLOAT,
		datatype.TINYTEXT, datatype.TEXT, datatype.MEDIUMTEXT, datatype.LONGTEXT:
		return 0
	// ==================================================
	default:
		return size
	}
}

// Protobuf 中的變數類型，轉為 SQL 中的變數類型；若為原生 SQL 變數類型，則無須修改
func (s *maria) ProtoTypeOf(kind string) string {
	switch kind {
	case "INT32":
		return datatype.INT
	case "INT64":
		return datatype.BIGINT
	case "BOOL":
		return datatype.TINYINT
	case "STRING":
		fallthrough
	case datatype.MESSAGE:
		fallthrough
	case datatype.MAP:
		return datatype.VARCHAR
	// 原生 SQL 變數類型，無須修改(EX: INT, TIMESTAMP)
	default:
		return kind
	}
}

func (s *maria) DbToProto(kind string) string {
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
func (s *maria) IsTableExistsStmt(tableName string) string {
	return fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE TABLE_NAME = '%s';", tableName)
}

// 是否為數值類型
func (s *maria) IsSortable(kind string) bool {
	return s.KindMap["STRING"].Contains(strings.ToUpper(kind))
}

// 判斷變數類型(integer, float, text, ...)
func (s *maria) GetKind(kind string) string {
	for k, v := range s.KindMap {
		if v.Contains(kind) {
			return k
		}
	}
	return "Unknown"
}
