package gdo

import (
	"reflect"

	"github.com/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// Insert
////////////////////////////////////////////////////////////////////////////////////////////////////
// 添加一筆數據(最終可同時添加多筆數據)
// 呼叫此函式者，須確保 datas 中的欄位都存在表格中
func (t *Table) Insert(data []any, ptrToDb func(reflect.Value, bool) string) error {
	var i int32
	length := t.GetColumnNumber()
	insertDatas := []string{}
	for i = 0; i < length; i++ {
		insertDatas = append(insertDatas, ValueToDb(reflect.ValueOf(data[i]), t.useAntiInjection, ptrToDb))
	}

	t.InsertStmt.Insert(insertDatas)
	return nil
}

func (t *Table) InsertRawData(data []string) {
	t.InsertStmt.Insert(data)
}

// 取得緩存數量
func (t *Table) GetInsertBufferNumber() int32 {
	return t.InsertStmt.GetBufferNumber()
}

func (t *Table) BuildInsertStmt() (string, error) {
	sql, err := t.InsertStmt.ToStmt()
	t.InsertStmt.Release()
	if err != nil {
		return "", errors.Wrap(err, "生成 InsertStmt 時發生錯誤")
	}
	return sql, nil
}
