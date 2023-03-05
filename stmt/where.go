package stmt

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type WhereStmt struct {
	sql          string
	useBrackets  bool
	notCondition bool
	ands         []*WhereStmt
	ors          []*WhereStmt
}

func WS() *WhereStmt {
	s := &WhereStmt{
		sql:          "",
		useBrackets:  false,
		notCondition: false,
		ands:         []*WhereStmt{},
		ors:          []*WhereStmt{},
	}
	return s
}

func (s *WhereStmt) SetBrackets() *WhereStmt {
	s.useBrackets = true
	return s
}

func (s *WhereStmt) SetNotCondition() *WhereStmt {
	s.notCondition = true
	return s
}

func (s *WhereStmt) AddAndCondtion(c *WhereStmt) *WhereStmt {
	s.ands = append(s.ands, c)
	return s
}

func (s *WhereStmt) ClearAndCondtion() {
	s.ands = s.ands[:0]
}

func (s *WhereStmt) AddOrCondtion(c *WhereStmt) *WhereStmt {
	s.ors = append(s.ors, c)
	return s
}

func (s *WhereStmt) ClearOrCondtion() {
	s.ors = s.ors[:0]
}

func (s *WhereStmt) Gt(key string, value string) *WhereStmt {
	s.sql = fmt.Sprintf("`%s` > %s", key, value)
	return s
}

func (s *WhereStmt) Ge(key string, value string) *WhereStmt {
	s.sql = fmt.Sprintf("`%s` >= %s", key, value)
	return s
}

func (s *WhereStmt) Eq(key string, value string) *WhereStmt {
	s.sql = fmt.Sprintf("`%s` = %s", key, value)
	return s
}

func (s *WhereStmt) Ne(key string, value string) *WhereStmt {
	s.sql = fmt.Sprintf("`%s` != %s", key, value)
	return s
}

func (s *WhereStmt) Le(key string, value string) *WhereStmt {
	s.sql = fmt.Sprintf("`%s` <= %s", key, value)
	return s
}

func (s *WhereStmt) Lt(key string, value string) *WhereStmt {
	s.sql = fmt.Sprintf("`%s` < %s", key, value)
	return s
}

func (s *WhereStmt) Like(key string, format string) *WhereStmt {
	s.sql = fmt.Sprintf("`%s` LIKE '%s'", key, format)
	return s
}

func (s *WhereStmt) Between(key string, value1 string, value2 string) *WhereStmt {
	s.sql = fmt.Sprintf("`%s` BETWEEN %s AND %s", key, value1, value2)
	return s
}

func (s *WhereStmt) In(key string, values ...string) *WhereStmt {
	s.sql = fmt.Sprintf("`%s` IN (%s)", key, strings.Join(values, ", "))
	return s
}

func (s *WhereStmt) CheckNull(key string, isNull bool) *WhereStmt {
	if isNull {
		s.sql = fmt.Sprintf("`%s` IS NULL", key)
	} else {
		s.sql = fmt.Sprintf("`%s` IS NOT NULL", key)
	}
	return s
}

func (s *WhereStmt) IsEmpty() bool {
	return s.sql == "" && len(s.ands) == 0 && len(s.ors) == 0
}

func (s *WhereStmt) Release() {
	s.sql = ""
	s.useBrackets = false
	s.notCondition = false
	s.ands = s.ands[:0]
	s.ors = s.ors[:0]
}

func (s *WhereStmt) ToStmt() (string, error) {
	if s.sql != "" {
		// 輸出前處理
		s.postStmtFormat()
		return s.sql, nil
	}

	if len(s.ands) != 0 {
		ands := []string{}
		var stmt string
		var err error

		for _, and := range s.ands {
			stmt, err = and.ToStmt()

			if err != nil {
				return "", errors.Wrap(err, "And condition has error.")
			}

			ands = append(ands, stmt)
		}

		s.sql = strings.Join(ands, " AND ")

		// 輸出前處理
		s.postStmtFormat()

		return s.sql, nil
	}

	if len(s.ors) != 0 {
		ors := []string{}
		var stmt string
		var err error

		for _, or := range s.ors {
			stmt, err = or.ToStmt()

			if err != nil {
				return "", errors.Wrap(err, "Or condition has error.")
			}

			ors = append(ors, stmt)
		}

		s.sql = strings.Join(ors, " OR ")

		// 輸出前處理
		s.postStmtFormat()

		return s.sql, nil
	}

	return "", nil
}

// 輸出前處理，檢查 括弧 與 NOT 語法
func (s *WhereStmt) postStmtFormat() {
	if s.useBrackets {
		s.sql = fmt.Sprintf("(%s)", s.sql)
	} else if s.notCondition {
		s.sql = fmt.Sprintf("NOT (%s)", s.sql)
	}
}
