package cls

import (
	"reflect"

	"github.com/j32u4ukh/gorm/gdo"
	"github.com/j32u4ukh/gorm/stmt"
)

func (t *StructTable) Update(data any, where *gdo.WhereStmt) (string, error) {
	var rv, field reflect.Value
	var column *stmt.Column
	var i int
	rv = reflect.ValueOf(data)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	// 遍歷每一欄位
	for i = 0; i < t.numFiled; i++ {
		field = rv.Field(i)
		// fmt.Printf("(t *ProtoTable) Update | field: %+v\n", field)
		column = t.Table.GetColumn(int32(i))

		if column.IgnoreThis {
			continue
		}

		switch column.Default {
		// 有值也不更新
		// timestamp 類型可透過設置 OnUpdate 來更新時間戳
		case "current_timestamp()", "AI":
			continue

		default:
			t.Table.UpdateRawData(column.Name, gdo.ValueToDb(field, t.useAntiInjection, StructToDb))
		}
	}

	if where != nil {
		t.Table.SetUpdateCondition(where)
	}

	sql, err := t.BuildUpdateStmt()
	return sql, err
}
