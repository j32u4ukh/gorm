package datatype

const (
	// Integer
	TINYINT   string = "TINYINT"
	SMALLINT  string = "SMALLINT"
	MEDIUMINT string = "MEDIUMINT"
	INT       string = "INT"
	BIGINT    string = "BIGINT"

	// Float
	FLOAT   string = "FLOAT"
	DOUBLE  string = "DOUBLE"
	DEMICAL string = "DEMICAL"

	// Text
	VARCHAR    string = "VARCHAR"
	CHAR       string = "CHAR"
	TINYTEXT   string = "TINYTEXT"
	TEXT       string = "TEXT"
	MEDIUMTEXT string = "MEDIUMTEXT"
	LONGTEXT   string = "LONGTEXT"

	// Time
	DATE      string = "DATE"
	TIME      string = "TIME"
	YEAR      string = "YEAR"
	DATETIME  string = "DATETIME"
	TIMESTAMP string = "TIMESTAMP"

	//////////////////////////////////////////////////
	// 非 SQL 類型 (Protobuf)
	//////////////////////////////////////////////////
	MAP      string = "MAP"
	MESSAGE  string = "MESSAGE"
	INT32    string = "INT32"
	INT64    string = "INT64"
	UINT32   string = "UINT32"
	UINT64   string = "UINT64"
	SINT32   string = "SINT32"
	SINT64   string = "SINT64"
	FIXED32  string = "FIXED32"
	FIXED64  string = "FIXED64"
	SFIXED32 string = "SFIXED32"
	SFIXED64 string = "SFIXED64"
	BOOL     string = "BOOL"
	STRING   string = "STRING"
	BYTES    string = "BYTES"
)

// 若為 MAP 或 proto 原生變數，則直接返回，其他則歸類為 MESSAGE
func GetOriginType(Type string) string {
	switch Type {
	case TINYINT, SMALLINT, MEDIUMINT, INT, BIGINT,
		FLOAT, DOUBLE, DEMICAL,
		VARCHAR, CHAR, TINYTEXT, TEXT, MEDIUMTEXT, LONGTEXT,
		DATE, TIME, YEAR, DATETIME, TIMESTAMP,
		MAP, MESSAGE,
		INT32, INT64, UINT32, UINT64, SINT32, SINT64,
		FIXED32, FIXED64, SFIXED32, SFIXED64,
		BOOL, STRING, "":
		return Type
	default:
		return MESSAGE
	}
}
