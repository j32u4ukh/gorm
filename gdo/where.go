package gdo

import (
	"reflect"

	"github.com/j32u4ukh/gorm/stmt"
)

type WhereStmt struct {
	*stmt.WhereStmt
	op     string
	values []any
	// 是否對 SQL injection 做處理
	useAntiInjection bool
	ands             []*WhereStmt
	ors              []*WhereStmt
	// 將變數反射為 SQL 數值的函式
	valueToDbFunc ValueToDbFunc
}

func WS() *WhereStmt {
	s := &WhereStmt{
		WhereStmt:        stmt.WS(),
		op:               "",
		values:           []any{},
		useAntiInjection: false,
		ands:             []*WhereStmt{},
		ors:              []*WhereStmt{},
		valueToDbFunc:    ValueToDb,
	}
	return s
}

func (s *WhereStmt) SetBrackets() *WhereStmt {
	s.WhereStmt.SetBrackets()
	return s
}

func (s *WhereStmt) SetNotCondition() *WhereStmt {
	s.WhereStmt.SetNotCondition()
	return s
}

func (s *WhereStmt) UseAntiInjection() *WhereStmt {
	s.useAntiInjection = true
	return s
}

func (s *WhereStmt) AddAndCondtion(c *WhereStmt) *WhereStmt {
	s.ands = append(s.ands, c)
	return s
}

func (s *WhereStmt) AddOrCondtion(c *WhereStmt) *WhereStmt {
	s.ors = append(s.ors, c)
	return s
}

func (s *WhereStmt) Gt(key string, value any) *WhereStmt {
	s.op = "Gt"
	s.values = []any{key, value}
	return s
}

func (s *WhereStmt) Ge(key string, value any) *WhereStmt {
	s.op = "Ge"
	s.values = []any{key, value}
	return s
}

func (s *WhereStmt) Eq(key string, value any) *WhereStmt {
	s.op = "Eq"
	s.values = []any{key, value}
	return s
}

func (s *WhereStmt) Ne(key string, value any) *WhereStmt {
	s.op = "Ne"
	s.values = []any{key, value}
	return s
}

func (s *WhereStmt) Le(key string, value any) *WhereStmt {
	s.op = "Le"
	s.values = []any{key, value}
	return s
}

func (s *WhereStmt) Lt(key string, value any) *WhereStmt {
	s.op = "Lt"
	s.values = []any{key, value}
	return s
}

func (s *WhereStmt) Like(key string, format any) *WhereStmt {
	s.op = "Like"
	s.values = []any{key, format}
	return s
}

func (s *WhereStmt) Between(key string, value1 any, value2 any) *WhereStmt {
	s.op = "Between"
	s.values = []any{key, value1, value2}
	return s
}

func (s *WhereStmt) In(key string, values ...any) *WhereStmt {
	s.op = "In"
	s.values = []any{key}
	s.values = append(s.values, values...)
	return s
}

func (s *WhereStmt) CheckNull(key string, isNull bool) *WhereStmt {
	s.op = "CheckNull"
	s.values = []any{key, isNull}
	return s
}

func (s *WhereStmt) IsEmpty() bool {
	return s.op == "" && len(s.ands) == 0 && len(s.ors) == 0
}

func (s *WhereStmt) Release() {
	s.WhereStmt.Release()
	s.op = ""
	s.values = s.values[:0]
	s.ands = s.ands[:0]
	s.ors = s.ors[:0]
}

func (s *WhereStmt) ToStmtWhere() *stmt.WhereStmt {
	// s.WhereStmt 還包含 useBrackets, notCondition 等其他設置，因此使用 s.WhereStmt 而非新建一個
	where := s.WhereStmt

	if s.op != "" {
		key := s.values[0].(string)
		if s.useAntiInjection {
			key = stmt.AntiInjectionString(key)
		}
		switch s.op {
		case "Gt":
			s.WhereStmt.Gt(key, s.valueToDbFunc(reflect.ValueOf(s.values[1]), s.useAntiInjection, nil))
		case "Ge":
			s.WhereStmt.Ge(key, s.valueToDbFunc(reflect.ValueOf(s.values[1]), s.useAntiInjection, nil))
		case "Eq":
			s.WhereStmt.Eq(key, s.valueToDbFunc(reflect.ValueOf(s.values[1]), s.useAntiInjection, nil))
		case "Ne":
			s.WhereStmt.Ne(key, s.valueToDbFunc(reflect.ValueOf(s.values[1]), s.useAntiInjection, nil))
		case "Le":
			s.WhereStmt.Le(key, s.valueToDbFunc(reflect.ValueOf(s.values[1]), s.useAntiInjection, nil))
		case "Lt":
			s.WhereStmt.Lt(key, s.valueToDbFunc(reflect.ValueOf(s.values[1]), s.useAntiInjection, nil))
		case "Like":
			// TODO: LIKE 的檢查應該和一般參數不同，應使用其他方式檢查
			s.WhereStmt.Like(key, s.valueToDbFunc(reflect.ValueOf(s.values[1]), s.useAntiInjection, nil))
		case "Between":
			s.WhereStmt.Between(key,
				s.valueToDbFunc(reflect.ValueOf(s.values[1]), s.useAntiInjection, nil),
				s.valueToDbFunc(reflect.ValueOf(s.values[2]), s.useAntiInjection, nil),
			)
		case "In":
			values := []string{}
			for _, value := range s.values[1:] {
				values = append(values, s.valueToDbFunc(reflect.ValueOf(value), s.useAntiInjection, nil))
			}
			s.WhereStmt.In(key, values...)
		case "CheckNull":
			s.WhereStmt.CheckNull(key, s.values[1].(bool))
		}
		return where
	}

	if len(s.ands) != 0 {
		where.ClearAndCondtion()
		var and *WhereStmt
		for _, and = range s.ands {
			if s.useAntiInjection {
				and.UseAntiInjection()
			}
			where.AddAndCondtion(and.ToStmtWhere())
		}
		return where
	}

	if len(s.ors) != 0 {
		where.ClearOrCondtion()
		var or *WhereStmt
		for _, or = range s.ors {
			if s.useAntiInjection {
				or.UseAntiInjection()
			}
			where.AddOrCondtion(or.ToStmtWhere())
		}
		return where
	}

	return where
}

func (s *WhereStmt) ToStmt() (string, error) {
	where := s.ToStmtWhere()
	return where.ToStmt()
}
