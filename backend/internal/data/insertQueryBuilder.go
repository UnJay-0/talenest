package data

import (
	"errors"
	"fmt"
	"strings"
)

type InsertQueryBuilder struct {
	query      strings.Builder
	ignore     bool
	replace    bool
	columns    []Column
	valueLists [][]Value
	tableName  string
}

func NewInsertQueryBuilder(tableName string) *InsertQueryBuilder {
	return &InsertQueryBuilder{
		tableName: tableName,
	}
}

func (builder *InsertQueryBuilder) SetIgnore() error {
	if builder.replace {
		return errors.New("You can't set ignore and replace at the same time")
	}
	builder.ignore = true
	return nil
}

func (builder *InsertQueryBuilder) UnsetIgnore() {
	builder.ignore = false
}

func (builder *InsertQueryBuilder) SetReplace() error {
	if builder.ignore {
		return errors.New("You can't set ignore and replace at the same time")
	}
	builder.replace = true
	return nil
}

func (builder *InsertQueryBuilder) UnsetReplace() {
	builder.replace = false
}

func (builder *InsertQueryBuilder) SetColumns(cols []Column) error {
	if builder.valueLists != nil && len(builder.valueLists) != len(cols) {
		return errors.New("The number of columns must be equal to the number of values")
	}
	builder.columns = cols
	return nil
}

func (builder *InsertQueryBuilder) AddColumn(col string) error {
	if builder.valueLists != nil && len(builder.valueLists) != (len(builder.columns)+1) {
		return errors.New("The number of columns must be equal to the number of values")
	}

	builder.columns = append(builder.columns, Column{colname: col, alias: ""})
	return nil
}

func (builder *InsertQueryBuilder) SetValues(values []Value) error {
	if builder.columns != nil && len(builder.columns) != len(values) {
		return errors.New("The number of columns must be equal to the number of values")
	}
	if len(builder.valueLists) > 0 && len(values) != len(builder.valueLists[len(builder.valueLists)-1]) {
		return errors.New("The number of values must correspond to the already inserted")
	}
	builder.valueLists = append(builder.valueLists, values)
	return nil
}

func (builder *InsertQueryBuilder) AddValue(value Value) error {
	if builder.columns != nil && len(builder.columns) != (len(builder.valueLists)+1) {
		return errors.New("The number of columns must be equal to the number of values")
	}
	if len(builder.valueLists) == 0 {
		builder.valueLists = append(builder.valueLists, []Value{})
	}
	builder.valueLists[len(builder.valueLists)-1] = append(
		builder.valueLists[len(builder.valueLists)-1], value)
	return nil
}

func buildValuesString(values []Value) string {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(fmt.Sprintf("%s", values[0]))

	for i := 1; i < len(values); i++ {
		builder.WriteString(fmt.Sprintf(", %s", values[i]))
	}

	builder.WriteString(")")
	return builder.String()
}

func (builder *InsertQueryBuilder) Build() (string, error) {
	if builder.valueLists == nil || len(builder.valueLists) == 0 {
		return "", errors.New("No value provided for the insert")
	}

	builder.query.WriteString("INSERT ")
	if builder.replace {
		builder.query.WriteString("OR REPLACE ")
	}
	if builder.ignore {
		builder.query.WriteString("OR IGNORE ")
	}

	builder.query.WriteString("INTO " + builder.tableName + " ")

	if builder.columns != nil && len(builder.columns) > 0 {
		builder.query.WriteString("(" + builder.columns[0].String())
		for i := 1; i < len(builder.columns); i++ {
			builder.query.WriteString(", " + builder.columns[i].String())
		}
		builder.query.WriteString(") ")
	}

	builder.query.WriteString("VALUES ")
	builder.query.WriteString(buildValuesString(builder.valueLists[0]))
	for i := 1; i < len(builder.valueLists); i++ {
		builder.query.WriteString(", " + buildValuesString(builder.valueLists[i]))
	}

	return builder.query.String(), nil
}

func (builder *InsertQueryBuilder) String() string {
	query, err := builder.Build()
	if err != nil {
		// TODO: Handle error
	}
	return query
}
