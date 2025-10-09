package data

import (
	"errors"
	"fmt"
	"strings"
)

type Column struct {
	colname string
	alias   string
}

func NewColumn(colName string, alias string) (*Column, error) {
	if colName == "" {
		return nil, errors.New("column names can't be empty while the alias is filled")
	}

	return &Column{
		colname: colName,
		alias:   alias,
	}, nil
}

func (column *Column) String() string {
	if column.alias == "" {
		return column.colname
	}
	return fmt.Sprintf("%s AS %s", column.colname, column.alias)
}

func (column *Column) GetColumnName() string {
	if column.alias != "" {
		return column.alias
	}
	return column.colname
}

func ConvertToColumns(columns []string) []Column {
	var result []Column
	for _, col := range columns {
		column, _ := NewColumn(col, "")
		result = append(result, *column)
	}
	return result
}

func ColumnsString(columns []Column) string {
	if len(columns) == 0 {
		return "*"
	}
	var builder strings.Builder

	builder.WriteString(columns[0].String())
	for i := 1; i < len(columns); i++ {
		builder.WriteString(", ")
		builder.WriteString(columns[i].String())
	}

	return builder.String()
}

type whereItem struct {
	tableName     string
	col           Column
	operator      string
	value         Value
	logicOperator string
}

func newWhereItem(tableName string, col Column, operator string, value Value, logicOperator string) *whereItem {
	// TODO: check valid operator
	// TODO: check valid logicOperator
	return &whereItem{
		tableName:     tableName,
		col:           col,
		operator:      operator,
		value:         value,
		logicOperator: logicOperator,
	}
}

func (item *whereItem) String() string {
	return fmt.Sprintf("%v.%v %v %v", item.tableName, item.col.String(), item.operator, item.value)
}

func whereString(whereItems []whereItem) string {
	if len(whereItems) == 0 {
		return ""
	}
	var builder strings.Builder

	builder.WriteString(whereItems[0].String())
	for i := 1; i < len(whereItems); i++ {
		builder.WriteString(whereItems[i].String() + " ")
		builder.WriteString(whereItems[i].logicOperator)

	}
	return builder.String()
}

type joinItem struct {
	sourceTable string
	sourceField string
	targetTable string
	targetField string
	joinType    string
}

func newJoinItem(sourceTable, sourceField, targetTable, targetField, joinType string) *joinItem {
	// TODO: Check joinType

	return &joinItem{
		sourceTable: sourceTable,
		sourceField: sourceField,
		targetField: targetField,
		targetTable: targetTable,
		joinType:    joinType,
	}
}

func (item *joinItem) String() string {
	return fmt.Sprintf("%s JOIN %s ON %s.%s = %s.%s ",
		item.joinType,
		item.targetTable,
		item.sourceTable,
		item.sourceField,
		item.targetTable,
		item.targetField)

}

type orderByItem struct {
	cols    []Column
	orderBy string
}

func newOrderByItem(cols []Column, orderBy string) *orderByItem {
	// TODO: check cols
	return &orderByItem{
		cols:    cols,
		orderBy: orderBy,
	}
}

func (item *orderByItem) String() string {
	builder := strings.Builder{}
	builder.WriteString("ORDER BY ")
	for i, col := range item.cols {
		builder.WriteString(col.GetColumnName())
		if i < len(item.cols)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString(" ")
	builder.WriteString(item.orderBy)
	return builder.String()
}

type SelectQueryBuilder struct {
	query     strings.Builder
	columns   []Column
	tableName string
	distinct  bool
	whereList []whereItem
	orderItem *orderByItem
	joinItems []joinItem
}

func NewSelectQueryBuilder(tableName string) *SelectQueryBuilder {
	return &SelectQueryBuilder{
		tableName: tableName,
	}
}

func (builder *SelectQueryBuilder) SetDistinct() {
	builder.distinct = true
}

func (builder *SelectQueryBuilder) SetColumns(columns []Column) {
	if len(columns) == 0 {
		col, _ := NewColumn("*", "")
		builder.columns = append(builder.columns, *col)
	}
	builder.columns = columns
}

func (builder *SelectQueryBuilder) SetWhere(tableName string, col Column, operator string, value Value, logicOperator string) {
	item := newWhereItem(tableName, col, operator, value, logicOperator)
	builder.whereList = append(builder.whereList, *item)
}

func (builder *SelectQueryBuilder) SetJoin(
	sourceTable string, sourceField Column, targetTable string, targetField Column, joinType string) {
	item := newJoinItem(sourceTable, sourceField.GetColumnName(), targetTable, targetField.GetColumnName(), joinType)
	builder.joinItems = append(builder.joinItems, *item)
}

func (builder *SelectQueryBuilder) OrderBy(columns []Column, orderBy string) {
	builder.orderItem = newOrderByItem(columns, orderBy)
}

func (builder *SelectQueryBuilder) Build() string {
	builder.query.WriteString("SELECT ")
	if builder.distinct {
		// Adding DISTINCT
		builder.query.WriteString("DISTINCT ")
	}

	// Adding columns
	builder.query.WriteString(ColumnsString(builder.columns))
	builder.query.WriteRune(' ')

	// Adding FROM
	builder.query.WriteString(fmt.Sprintf("FROM %s ", builder.tableName))

	// Adding joins
	if builder.joinItems != nil || len(builder.joinItems) > 0 {
		for _, join := range builder.joinItems {
			builder.query.WriteString(join.String())
		}
	}

	// Adding where
	if builder.whereList != nil || len(builder.whereList) > 0 {
		builder.query.WriteString("WHERE ")
		builder.query.WriteString(whereString(builder.whereList))
	}

	// Adding order by
	if builder.orderItem != nil {
		builder.query.WriteString(builder.orderItem.String())
	}

	builder.query.WriteRune(';')
	return builder.query.String()
}
