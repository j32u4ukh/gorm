package stmt

import (
	"fmt"
	"strings"

	"github.com/j32u4ukh/cntr"
)

// CREATE DATABASE `PVP` /*!40100 COLLATE 'utf8mb4_bin' */
type CreateStmt struct {
	DbName     string
	TableName  string
	TableParam *TableParam
	Columns    []*Column
	Engine     string
	Collate    string
}

func NewCreateStmt(name string, tableParam *TableParam, columnParams []*ColumnParam, engine string, collate string) *CreateStmt {
	s := &CreateStmt{
		DbName:     "",
		TableName:  name,
		TableParam: tableParam,
		Columns:    []*Column{},
		Engine:     engine,
		Collate:    collate,
	}

	if columnParams != nil {
		var column *Column
		for _, param := range columnParams {
			fmt.Printf("param: %+v\n", param)
			column = NewColumn(param)
			column.SetCollate(collate)
			s.AddColumn(column)
		}
	}
	return s
}

func (s *CreateStmt) SetDbName(dbName string) *CreateStmt {
	s.DbName = dbName
	return s
}

func (s *CreateStmt) SetEngine(engine string) {
	s.Engine = engine
}

func (s *CreateStmt) SetCollate(collate string) {
	s.Collate = collate
}

func (s *CreateStmt) GetTableParam() *TableParam {
	return s.TableParam
}

func (s *CreateStmt) AddColumn(column *Column) *CreateStmt {
	s.Columns = append(s.Columns, column)
	if column.IsPrimaryKey {
		s.TableParam.AddPrimaryKey(column.Name, column.Algo)
	}
	return s
}

func (s *CreateStmt) ToStmt() (string, error) {
	// CREATE TABLE `DESK` (
	// 		"`id` INT NOT NULL",
	// 		"`text` VARCHAR(45) NULL",
	// 		"`timestamp` BIGINT(13) NULL",
	// 		"`flag` TINYINT NULL",
	// 		"PRIMARY KEY (`id`)",
	// )ENGINE=InnoDB COLLATE=utf8mb4_bin
	stmts := []string{}
	var column *Column

	for _, column = range s.Columns {
		if column.IgnoreThis {
			continue
		}
		stmts = append(stmts, column.ToStmt())
	}

	// "PRIMARY KEY (`key_column_name`)"
	pks := []string{}
	for _, key := range s.TableParam.Primarys.Elements {
		// s.TableParam.Primarys.Elements[i] = fmt.Sprintf("`%s`", key)
		pks = append(pks, fmt.Sprintf("`%s`", key))
	}

	stmts = append(stmts, fmt.Sprintf("PRIMARY KEY (%s) USING %s",
		strings.Join(pks, ", "),
		s.TableParam.IndexType["PRIMARY"],
	))

	if s.TableParam != nil {
		var kind, indexName, indexType string
		var cols []string

		// 0: kind string, 1: indexName string, 2: indexType string, 3: cols *array.Array[string]
		it := s.TableParam.IterIndexMap()

		for it.HasNext() {
			data := it.Next().([]any)
			kind = data[0].(string)
			indexName = data[1].(string)
			indexType = data[2].(string)
			cols = data[3].(*cntr.Array[string]).Elements

			for i, col := range cols {
				cols[i] = fmt.Sprintf("`%s`", col)
			}

			indexStmt := fmt.Sprintf("INDEX `%s` (%s) USING %s",
				indexName, strings.Join(cols, ", "), indexType)

			if kind == "UNIQUE" {
				indexStmt = fmt.Sprintf("UNIQUE %s", indexStmt)
			}

			stmts = append(stmts, indexStmt)
		}
	}

	var tableName string

	if s.DbName != "" {
		tableName = fmt.Sprintf("`%s`.`%s`", s.DbName, s.TableName)
	} else {
		tableName = fmt.Sprintf("`%s`", s.TableName)
	}

	sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (%s) ENGINE = %s COLLATE = '%s';`,
		tableName,
		strings.Join(stmts, ", "),
		s.Engine,
		s.Collate,
	)

	return sql, nil
}
