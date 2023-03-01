package stmt

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// UpdateStmt
////////////////////////////////////////////////////////////////////////////////////////////////////
type UpdateStmt struct {
	DbName    string
	TableName string
	datas     []string
	Where     *WhereStmt
}

func NewUpdateStmt(tableName string) *UpdateStmt {
	s := &UpdateStmt{
		DbName:    "",
		TableName: tableName,
		// pairs:     NewKeyValueSlice(),
		datas: []string{},
		Where: nil,
	}
	return s
}

func (s *UpdateStmt) SetDbName(dbName string) *UpdateStmt {
	s.DbName = dbName
	return s
}

// 要求外部依據固定傳入順序，避免還要進行排序
func (s *UpdateStmt) Update(key string, value string) *UpdateStmt {
	// pair := NewKeyValuePair(key, value)
	// fmt.Printf("(s *UpdateStmt) Update | %+v\n", pair)
	// s.pairs.AddPair(pair)
	s.datas = append(s.datas, fmt.Sprintf("`%s` = %s", key, value))
	return s
}

func (s *UpdateStmt) SetCondition(where *WhereStmt) *UpdateStmt {
	s.Where = where
	return s
}

func (s *UpdateStmt) Release() {
	s.datas = s.datas[:0]
	s.Where = nil
	// s.pairs.Release()
	// s.Where.Release()
}

/*
   修改非 primary key 欄位
   UPDATE `demo2`.`Desk` SET `item_id`='3' WHERE  `index`=1 AND `user_id`=2;

   修改 primary key 欄位
   UPDATE `demo2`.`Desk` SET `user_id`='5' WHERE  `index`=1 AND `user_id`=2;
*/
func (s *UpdateStmt) ToStmt() (string, error) {
	if len(s.datas) == 0 {
		return "", errors.New("Update data is empty.")
	}

	// if s.pairs.IsEmpty() {
	// 	return "", errors.New("Update data is empty.")
	// }

	if s.Where == nil {
		return "", errors.New("Where condition is nil.")
	}

	where, err := s.Where.ToStmt()

	if err != nil || where == "" {
		return "", errors.Wrapf(err, "Failed to create where statment.")
	}

	var tableName string

	if s.DbName != "" {
		tableName = fmt.Sprintf("`%s`.`%s`", s.DbName, s.TableName)
	} else {
		tableName = fmt.Sprintf("`%s`", s.TableName)
	}

	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s;", tableName, strings.Join(s.datas, ", "), where)
	return sql, nil
}

// NOTE: 批次更新仍保留，但目前仍有使用 SqlValue(value)，待之後有空再來修改
// ////////////////////////////////////////////////////////////////////////////////////////////////////
// // BatchUpdateStmt
// ////////////////////////////////////////////////////////////////////////////////////////////////////
// type BatchUpdateStmt struct {
// 	Name        string
// 	PrimaryKey  string
// 	PrimaryKeys []string
// 	// Key: column name, Value: SetStmt
// 	sets map[string]*SetStmt
// }

// // primaryKey: 多組數據時，根據此欄位來區分不同數據
// func NewBatchUpdateStmt(name string, primaryKey string) *BatchUpdateStmt {
// 	s := &BatchUpdateStmt{
// 		Name:        name,
// 		PrimaryKey:  primaryKey,
// 		PrimaryKeys: []string{},
// 		sets:        map[string]*SetStmt{},
// 	}
// 	return s
// }

// func (s *BatchUpdateStmt) Update(data map[string]any) *BatchUpdateStmt {
// 	key := SqlValue(data[s.PrimaryKey])
// 	s.PrimaryKeys = append(s.PrimaryKeys, key)
// 	var ok bool

// 	for col, value := range data {
// 		if col == s.PrimaryKey {
// 			continue
// 		}

// 		if _, ok = s.sets[col]; !ok {
// 			s.sets[col] = newSetStmt(s.PrimaryKey, col)
// 		}

// 		s.sets[col].AddData(key, value)
// 	}
// 	return s
// }

// // 取得緩存數量
// func (s *BatchUpdateStmt) GetBufferNumber() int {
// 	return len(s.PrimaryKeys)
// }

// func (s *BatchUpdateStmt) Release() {
// 	s.PrimaryKeys = []string{}

// 	for k := range s.sets {
// 		delete(s.sets, k)
// 	}
// }

// func (s *BatchUpdateStmt) ToStmt() (string, error) {
// 	sets := []string{}
// 	var stmt string
// 	var err error

// 	for _, set := range s.sets {
// 		stmt, err = set.toStmt()

// 		if err != nil {
// 			return "", errors.Wrap(err, "Failed to generate set statement.")
// 		}

// 		sets = append(sets, stmt)
// 	}

// 	setStmt := strings.Join(sets, ", ")
// 	pks := strings.Join(s.PrimaryKeys, ", ")
// 	sql := fmt.Sprintf("UPDATE %s SET %s WHERE `%s` IN (%s);", s.Name, setStmt, s.PrimaryKey, pks)
// 	return sql, nil
// }

// ////////////////////////////////////////////////////////////////////////////////////////////////////
// // SetStmt
// ////////////////////////////////////////////////////////////////////////////////////////////////////
// type SetStmt struct {
// 	key    string
// 	column string
// 	// key: value of primary column; value: value of target column
// 	data map[string]string
// }

// func newSetStmt(key string, column string) *SetStmt {
// 	s := &SetStmt{
// 		key:    key,
// 		column: column,
// 		data:   map[string]string{},
// 	}
// 	return s
// }

// func (s *SetStmt) AddData(key string, value any) {
// 	s.data[key] = SqlValue(value)
// 	// fmt.Printf("key: %s, value: %s\n", key, s.data[key])
// }

// func (s *SetStmt) toStmt() (string, error) {
// 	// SET [column] = CASE [primary_key]
// 	// WHEN m0.fileds.Get(0) THEN 'Insert1'
// 	// WHEN m1.fileds.Get(0) THEN 'Insert3'
// 	// END
// 	content := []string{}

// 	// 確保每次輸出順序相同
// 	keys := make([]string, 0, len(s.data))
// 	for k := range s.data {
// 		// fmt.Printf("(s *SetStmt) toStmt | k: %s\n", k)
// 		keys = append(keys, k)
// 	}
// 	sort.Strings(keys)
// 	// fmt.Printf("(s *SetStmt) toStmt | keys: %+v\n", keys)

// 	for _, key := range keys {
// 		// WHEN [value of primary column] THEN [value of target column]
// 		content = append(content, fmt.Sprintf("WHEN %s THEN %s", key, s.data[key]))
// 	}

// 	return fmt.Sprintf("`%s` = CASE `%s` %s END", s.column, s.key, strings.Join(content, " ")), nil
// }
