package database

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

var dbMap map[int32]*Database
var mu sync.Mutex

func init() {
	dbMap = map[int32]*Database{}
}

func Connect(idx int32, user string, password string, host string, port int, dbName string) (*Database, error) {
	if _, ok := dbMap[idx]; ok {
		return nil, errors.New(fmt.Sprintf("Database(%d) 已被使用", idx))
	}

	dm, err := NewDatabase(idx, user, password, host, port)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to new Database.")
	}

	err = dm.Use(dbName)

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to use database `%s`.", dbName)
	}

	mu.Lock()
	defer mu.Unlock()
	dbMap[idx] = dm
	return dm, nil
}

func Get(idx int32) *Database {
	if dm, ok := dbMap[idx]; ok {
		return dm
	}
	return nil
}

func Close(idx int32) error {
	dm := Get(idx)
	if dm == nil {
		return errors.New(fmt.Sprintf("不存在標號為 %d 的資料庫物件", idx))
	}
	mu.Lock()
	defer mu.Unlock()
	err := dm.close()
	if err != nil {
		return errors.Wrapf(err, "關閉號為 %d 的資料庫物件時發生錯誤", idx)
	}
	delete(dbMap, idx)
	return nil
}
