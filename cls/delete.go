package cls

import (
	"github.com/j32u4ukh/gorm/gdo"
	"github.com/pkg/errors"
)

func (t *StructTable) BuildDeleteStmt(where *gdo.WhereStmt) (string, error) {
	if where != nil {
		t.Table.SetDeleteCondition(where)
	}

	sql, err := t.Table.BuildDeleteStmt()

	if err != nil {
		return "", errors.Wrap(err, "Failed to build DeleteStmt.")
	}

	return sql, nil
}
