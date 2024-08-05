package database

import (
	"errors"
	"fmt"
)

const QUERY_NAME = "relations"

type Relation struct {
	ConstraintName    string `db:"constraint_name"`
	TableName         string `db:"table_name"`
	ColumnName        string `db:"column_name"`
	ForeignTableName  string `db:"foreign_table_name"`
	ForeignColumnName string `db:"foreign_column_name"`
}

func GetRelations(table string) ([]Relation, error) {
	var rels []Relation
	err := QuerySelect(QUERY_NAME, &rels, table)
	return rels, err
}

func LoadRelations(model Model) error {
	rels, err := GetRelations(model.TableName())
	if err != nil {
		return err
	}
	associations := model.Associations()
	for _, rel := range rels {
		for _, a := range associations {
			relModel, ok := a.(Model)
			if !ok {
				return errors.New("relation is not Model compatible")
			}
			if rel.TableName == relModel.TableName() {
				fmt.Println(rel, "^^^^^^^^^^^^^^^^^^^^^")
			}
		}
	}
	return nil
}
