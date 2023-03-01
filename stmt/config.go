package stmt

import (
	"encoding/json"

	"github.com/j32u4ukh/cntr"
	"github.com/pkg/errors"
)

type TableParamConfig struct {
	Uniques []IndexConfig `json:"unique"`
	Indexs  []IndexConfig `json:"index"`
}

type IndexConfig struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Columns []string `json:"columns"`
}

func NewTableParamConfig(data string) (*TableParamConfig, error) {
	tpc := &TableParamConfig{}
	err := json.Unmarshal([]byte(data), tpc)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal config.")
	}

	return tpc, nil
}

type ColumnParamConfig struct {
	Type string `json:"type"`
	Size int    `json:"size"`
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
	err := json.Unmarshal([]byte(data), cpc)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal config.")
	}

	return cpc, nil
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
