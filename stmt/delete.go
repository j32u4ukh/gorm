package stmt

import (
	"fmt"

	"github.com/pkg/errors"
)

type DeleteStmt struct {
	DbName string
	Name   string
	Where  *WhereStmt
}

func NewDeleteStmt(name string) *DeleteStmt {
	s := &DeleteStmt{
		DbName: "",
		Name:   name,
		Where:  &WhereStmt{},
	}
	return s
}

func (s *DeleteStmt) SetDbName(dbName string) *DeleteStmt {
	s.DbName = dbName
	return s
}

func (s *DeleteStmt) SetCondition(where *WhereStmt) *DeleteStmt {
	s.Where = where
	return s
}

func (s *DeleteStmt) Release() *DeleteStmt {
	s.Where.Release()
	return s
}

func (s *DeleteStmt) ToStmt() (string, error) {
	// DELETE FROM `demo`.`dartnotification` WHERE (`id` = '1665154635370141');
	// 一次刪除某資料表中所有的資料：DELETE FROM table_name; | DELETE * FROM table_name;
	where, err := s.Where.ToStmt()

	if err != nil {
		return "", errors.Wrap(err, "Failed to generate where statement.")
	}

	// 檢查 where 是否為空，否則就會形成 刪除資料表中所有的資料 的語法
	if where != "" {
		where = fmt.Sprintf(" WHERE %s", where)
	}

	var tableName string

	if s.DbName != "" {
		tableName = fmt.Sprintf("`%s`.`%s`", s.DbName, s.Name)
	} else {
		tableName = fmt.Sprintf("`%s`", s.Name)
	}

	sql := fmt.Sprintf("DELETE FROM %s%s;", tableName, where)
	return sql, nil
}
