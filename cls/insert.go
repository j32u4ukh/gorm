package cls

import (
	"reflect"

	"github.com/j32u4ukh/gorm/gdo"
	"github.com/j32u4ukh/gorm/stmt"
)

func (t *StructTable) Insert(datas []any) error {
	var data any
	var rv, field reflect.Value
	var column *stmt.Column
	var i, idx int

	for _, data = range datas {
		rv = reflect.ValueOf(data).Elem()
		data := []string{}

		for i = 0; i < t.numFiled; i++ {
			column = t.Table.GetColumn(int32(idx))

			if column.IgnoreThis {
				continue
			}

			switch column.Default {
			// 資料庫自動生成欄位
			case "current_timestamp()", "AI":
				data = append(data, "NULL")
			default:
				field = rv.FieldByIndex([]int{i})
				data = append(data, gdo.ValueToDb(field, t.useAntiInjection, StructToDb))
			}
		}

		// 將數據加入 insert 緩存(傳入數據由 ProtoMessage 生成，因此所有欄位一定都在，且依照欄位順序)
		t.Table.InsertRawData(data...)
	}
	return nil
}
