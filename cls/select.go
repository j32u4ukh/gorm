package cls

import (
	"reflect"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gorm/gdo"

	"github.com/pkg/errors"
)

// 根據 ProtoMessage 生成 SQL 的查詢語法，m 須包含 table 的 primary key
func (t *StructTable) BuildSelectStmt() (string, error) {
	sql, err := t.Table.BuildSelectStmt()

	if err != nil {
		return "", err
	}

	return sql, nil
}

func (t *StructTable) ParseSelectResults(datas *[]any, results [][]string) error {
	var i, length int32 = 0, int32(len(results))
	pm := (*datas)[0]
	t.InitByStruct(pm)
	var err error

	for i = 0; i < length; i++ {
		err = t.parseSelectResult((*datas)[i], results[i])
		if err != nil {
			return errors.Wrapf(err, "解析回傳數據時發生錯誤, result: %s", cntr.SliceToString(results[i]))
		}
	}
	return nil
}

func (t *StructTable) parseSelectResult(data any, result []string) error {
	rv := reflect.ValueOf(data).Elem()
	var filed reflect.Value

	for i, res := range result {
		if res == "" {
			continue
		}
		filed = rv.FieldByIndex([]int{i + 3})
		gdo.SetValue(filed, []byte(res), nil)
	}
	return nil
}

// 取得符合 WhereStmt 條件的數據筆數
func (t *StructTable) CountStmt(where *gdo.WhereStmt) string {
	t.CountMode()
	t.SetSelectCondition(where)
	sql, err := t.BuildSelectStmt()
	if err != nil {
		return ""
	}
	return sql
}
