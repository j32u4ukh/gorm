package cls

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/j32u4ukh/gorm/stmt"
	"github.com/j32u4ukh/gorm/stmt/dialect"
)

type ISqlStruct interface {
	ToStmt() string
}

func GetParams(data any, dial dialect.SQLDialect) (*stmt.TableParam, []*stmt.ColumnParam, error) {
	rv := reflect.ValueOf(data)
	var rt reflect.Type

	if rv.Kind() == reflect.Ptr {
		rt = rv.Elem().Type()
	} else {
		rt = reflect.TypeOf(data)
	}

	tableParam := stmt.NewTableParam()
	columnParams := []*stmt.ColumnParam{}
	var tpc *stmt.TableParamConfig
	var cpc *stmt.ColumnParamConfig

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		columnParam := stmt.NewColumnParam(i, field.Name, field.Type.Kind().String(), dial)
		config, ok := field.Tag.Lookup("gorm")

		if ok {
			tpc, cpc, _ = stmt.ParseConfig(field.Name, config)
			tableParam.LoadConfig(tpc)
			columnParam.LoadConfig(cpc)
		} else {
			fmt.Printf("No tag found for field %s\n", field.Name)
		}

		columnParam.Redefine()
		columnParams = append(columnParams, columnParam)
	}

	return tableParam, columnParams, nil
}

func StructToDb(field reflect.Value, useAntiInjection bool) string {
	m := field.Interface().(ISqlStruct)
	return m.ToStmt()
}

func GetTagName(tag string) (string, bool) {
	pairs := strings.Split(tag, ",")
	for _, pair := range pairs {
		k, v, ok := strings.Cut(pair, "=")
		if ok && k == "name" {
			return v, true
		}
	}
	return "", false
}
