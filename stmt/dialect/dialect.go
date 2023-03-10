package dialect

type SQLDialect byte

const (
	MARIA SQLDialect = iota
	MYSQL
	SQLITE3
)

var dialectsMap = map[SQLDialect]IDialect{}

// 處理不同資料庫在語法上的差異
type IDialect interface {
	// 變數類型，轉為 SQL 中的變數類型
	TypeOf(dataType string) string

	// 根據 dataType 、當前的 size 以及 DB 本身的限制，對數值大小再定義
	SizeOf(dataType string, size int32) int32

	GetDefault(dataType string) string

	// // TODO: 移除此函式，Dialect 本身不應與 Protobuf 有關聯或耦合
	// // Protobuf 中的變數類型，轉為 SQL 中的變數類型
	// ProtoTypeOf(string) string

	// // SQL 中的變數類型，轉為 Protobuf 中的變數類型
	// DbToProto(string) string

	// // 表格是否存在的 SQL 語法
	// IsTableExistsStmt(tableName string) string

	// 是否可適用排序演算法
	IsSortable(kind string) bool

	// 判斷變數類型(integer, float, text, ...)
	GetKind(kind string) string
}

// 註冊各資料庫語言的方言物件
func RegisterDialect(name SQLDialect, dialect IDialect) {
	dialectsMap[name] = dialect
}

// 取得資料庫語言的方言物件
func GetDialect(name SQLDialect) IDialect {
	return dialectsMap[name]
}
