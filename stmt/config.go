package stmt

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/j32u4ukh/cntr"
	"github.com/pkg/errors"
)

func ParseConfig(columnName string, config string) (*TableParamConfig, *ColumnParamConfig, error) {
	tpc, _ := NewTableParamConfig("")
	cpc, _ := NewColumnParamConfig("")
	settings := strings.Split(config, ";")
	settingMap := make(map[string]string)
	var kv []string
	var setting, key, value string
	for _, setting = range settings {
		kv = strings.Split(setting, "=")
		if len(kv) != 2 {
			continue
		}
		key = strings.ToUpper(strings.Trim(kv[0], " "))
		settingMap[key] = strings.Trim(kv[1], " ")
		cpc.setConfig(key, settingMap[key])
		if key == "NAME" {
			columnName = settingMap[key]
		}
	}
	for key, value = range settingMap {
		switch key {
		case "PK":
			tpc.Primarys.Append(columnName)
		case "UNIQUE":
			names := strings.Split(value, ",")
			for _, name := range names {
				tpc.addUnique(name, "default", columnName)
			}
		case "INDEX":
			names := strings.Split(value, ",")
			for _, name := range names {
				tpc.AddIndex(name, "default", columnName)
			}
		}
	}
	return tpc, cpc, nil
}

type IndexConfig struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Columns []string `json:"columns"`
}

type TableParamConfig struct {
	Primarys *cntr.Array[string] `json:"primary"`
	Uniques  []*IndexConfig      `json:"unique"`
	Indexs   []*IndexConfig      `json:"index"`
}

func (tpc *TableParamConfig) addUnique(name string, kind string, column ...string) {
	var config *IndexConfig
	var i int
	for i, config = range tpc.Uniques {
		if config.Name == name {
			config.Columns = append(config.Columns, column...)
			tpc.Uniques[i] = config
			return
		}
	}

	tpc.Uniques = append(tpc.Uniques, &IndexConfig{
		Name:    name,
		Type:    kind,
		Columns: column,
	})
}

func (tpc *TableParamConfig) AddIndex(name string, kind string, column ...string) {
	var config *IndexConfig
	var i int
	for i, config = range tpc.Indexs {
		if config.Name == name {
			config.Columns = append(config.Columns, column...)
			tpc.Indexs[i] = config
			return
		}
	}

	tpc.Indexs = append(tpc.Indexs, &IndexConfig{
		Name:    name,
		Type:    kind,
		Columns: column,
	})
}

func NewTableParamConfig(data string) (*TableParamConfig, error) {
	tpc := &TableParamConfig{
		Primarys: cntr.NewArray[string](),
		Uniques:  []*IndexConfig{},
		Indexs:   []*IndexConfig{},
	}
	if data != "" {
		err := json.Unmarshal([]byte(data), tpc)

		if err != nil {
			return nil, errors.Wrap(err, "Failed to unmarshal config.")
		}
	}
	return tpc, nil
}

type ColumnParamConfig struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Size   int    `json:"size"`
	// 定義 PrimaryKey 所使用的排序演算法(若不特別定義，可使用 Default)
	PrimaryKey string `json:"primary_key"`
	Unsigned   string `json:"unsigned"`
	CanNull    string `json:"can_null"`
	Default    string `json:"default"`
	Update     string `json:"update"`
	Comment    string `json:"comment"`
}

func NewColumnParamConfig(data string) (*ColumnParamConfig, error) {
	cpc := &ColumnParamConfig{}
	if data != "" {
		err := json.Unmarshal([]byte(data), cpc)

		if err != nil {
			return nil, errors.Wrap(err, "Failed to unmarshal config.")
		}
	}
	return cpc, nil
}

func (c *ColumnParamConfig) setConfig(key string, value string) {
	switch key {
	case "NUMBER":
		c.Number, _ = strconv.Atoi(value)
	case "NAME":
		c.Name = value
	case "TYPE":
		c.Type = value
	case "SIZE":
		c.Size, _ = strconv.Atoi(value)
	case "PK":
		c.PrimaryKey = value
	case "UNSIGNED":
		c.Unsigned = value
	case "DEFAULT":
		c.Default = value
	case "UPDATE":
		c.Update = value
	case "COMMENT":
		c.Comment = value
	}
}

func (c *ColumnParamConfig) merge(other *ColumnParamConfig) {
	c.Type = mergeString(c.Type, other.Type)
	c.Size = mergeInt(c.Size, other.Size)
	c.PrimaryKey = mergeString(c.PrimaryKey, other.PrimaryKey)
	c.Unsigned = mergeString(c.Unsigned, other.Unsigned)
	c.CanNull = mergeString(c.CanNull, other.CanNull)
	c.Default = mergeString(c.Default, other.Default)
	c.Update = mergeString(c.Update, other.Update)
	c.Comment = mergeString(c.Comment, other.Comment)
}

func mergeString(origin string, source string) string {
	if source != "" {
		return source
	} else {
		return origin
	}
}

func mergeInt[T cntr.Int | cntr.UInt](origin T, source T) T {
	if source != 0 {
		return source
	} else {
		return origin
	}
}
