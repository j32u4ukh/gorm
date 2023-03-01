package stmt

import (
	"fmt"

	"github.com/pkg/errors"
)

type SelectStmt struct {
	DbName    string
	TableName string
	// 要查詢的欄位名稱列表
	QueryColumns []string
	// 查詢的篩選機制
	Where     *WhereStmt
	OrderBy   string
	OrderType string
	Limit     int32
	Offset    int32
}

func NewSelectStmt(tableName string) *SelectStmt {
	s := &SelectStmt{
		DbName:       "",
		TableName:    tableName,
		QueryColumns: []string{},
		Where:        &WhereStmt{},
		OrderBy:      "",
		OrderType:    "ASC",
		Limit:        -1,
		Offset:       0,
	}
	return s
}

func (s *SelectStmt) SetDbName(name string) *SelectStmt {
	s.DbName = name
	return s
}

func (s *SelectStmt) Query(column string) *SelectStmt {
	s.QueryColumns = append(s.QueryColumns, column)
	return s
}

// 利用原本 Columns 的機制，直接形成 SELECT COUNT(*) 的 SQL 語法
func (s *SelectStmt) CountMode() *SelectStmt {
	s.QueryColumns = append(s.QueryColumns, "COUNT(*)")
	return s
}

func (s *SelectStmt) SetCondition(where *WhereStmt) *SelectStmt {
	s.Where = where
	return s
}

func (s *SelectStmt) SetOrderBy(column string) *SelectStmt {
	s.OrderBy = column
	return s
}

// 升序(ASC) | 降序(DESC)
func (s *SelectStmt) WhetherReverseOrder(reverse bool) *SelectStmt {
	if reverse {
		s.OrderType = "DESC"
	} else {
		s.OrderType = "ASC"
	}
	return s
}

func (s *SelectStmt) SetLimit(limit int32) *SelectStmt {
	s.Limit = limit
	return s
}

func (s *SelectStmt) SetOffset(offset int32) *SelectStmt {
	s.Offset = offset
	return s
}

func (s *SelectStmt) Release() {
	s.QueryColumns = []string{}
	s.Where.Release()
	s.OrderBy = ""
	s.OrderType = "ASC"
	s.Limit = -1
	s.Offset = 0
}

func (s *SelectStmt) ToStmt() (string, error) {
	formatColumns := FormatColumns(s.QueryColumns)
	where, err := s.Where.ToStmt()

	if err != nil {
		return "", errors.Wrap(err, "Failed to generate where statement.")
	}

	if where != "" {
		where = fmt.Sprintf(" WHERE %s", where)
	}

	var order string
	var limitOffset string

	if s.OrderBy == "" {
		order = ""
	} else {
		order = fmt.Sprintf(" ORDER BY `%s` %s", s.OrderBy, s.OrderType)
	}

	if s.Limit == -1 {
		limitOffset = ""
	} else {
		limitOffset = fmt.Sprintf(" LIMIT %d OFFSET %d", s.Limit, s.Offset)
	}

	var tableName string

	if s.DbName != "" {
		tableName = fmt.Sprintf("`%s`.`%s`", s.DbName, s.TableName)
	} else {
		tableName = fmt.Sprintf("`%s`", s.TableName)
	}

	sql := fmt.Sprintf("SELECT %s FROM %s%s%s%s;", formatColumns, tableName, where, order, limitOffset)
	return sql, nil
}
