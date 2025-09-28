package database

import (
	"strings"
)

type DeleteQueryBuilder struct {
	query     strings.Builder
	tableName string
	whereList []whereItem
}

func NewDeleteQueryBuilder(tableName string) *DeleteQueryBuilder {
	return &DeleteQueryBuilder{
		tableName: tableName,
	}
}

func (builder *DeleteQueryBuilder) SetWhere(tableName string, col Column, operator string, value Value, logicOperator string) {
	item := newWhereItem(tableName, col, operator, value, logicOperator)
	builder.whereList = append(builder.whereList, *item)
}

func (builder *DeleteQueryBuilder) Build() string {
	builder.query.WriteString("DELETE FROM ")
	builder.query.WriteString(builder.tableName)
	builder.query.WriteRune(' ')

	// Adding where
	if builder.whereList != nil || len(builder.whereList) > 0 {
		builder.query.WriteString("WHERE ")
		builder.query.WriteString(whereString(builder.whereList))
	}

	builder.query.WriteRune(';')
	return builder.query.String()
}
