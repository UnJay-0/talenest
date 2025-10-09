package data

import (
	"errors"
	"fmt"
)

type UpdateQueryBuilder struct {
	tableName string
	columns   []Column
	values    []Value
	whereList []whereItem
}

func NewUpdateQueryBuilder(tableName string) *UpdateQueryBuilder {
	return &UpdateQueryBuilder{
		tableName: tableName,
		columns:   []Column{},
		values:    []Value{},
		whereList: []whereItem{},
	}
}

func (builder *UpdateQueryBuilder) SetNewValues(columns []Column, values []Value) error {
	if len(columns) != len(values) {
		return errors.New("columns and values must have the same length")
	}

	builder.columns = columns
	builder.values = values
	return nil
}

func (builder *UpdateQueryBuilder) SetWhere(tableName string, col Column, operator string, value Value, logicOperator string) {
	item := newWhereItem(tableName, col, operator, value, logicOperator)
	builder.whereList = append(builder.whereList, *item)
}

func (builder *UpdateQueryBuilder) Build() string {
	query := fmt.Sprintf("UPDATE %s SET ", builder.tableName)
	for i, column := range builder.columns {
		query += fmt.Sprintf("%s = %s", column.GetColumnName(), builder.values[i].GetValueString())
		if i < len(builder.columns)-1 {
			query += ", "
		}
	}

	if len(builder.whereList) > 0 {
		query += " WHERE "
		query += whereString(builder.whereList)
	}
	query += ";"
	return query
}
