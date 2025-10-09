package data

const (
	CREATE_STATEMENT   = "CREATE"
	READ_STATEMENT     = "READ"
	READ_ALL_STATEMENT = "READ_ALL"
	UPDATE_STATEMENT   = "UPDATE"
	DELETE_STATEMENT   = "DELETE"
)

func CreateQuery(tableName string, columnNames []string) string {
	builder := NewInsertQueryBuilder(tableName)
	builder.SetColumns(ConvertToColumns(columnNames))
	builder.SetValues(GetTokens(len(columnNames), "?"))
	queryStr, _ := builder.Build()
	return queryStr
}

func ReadByIdQuery(tableName string, columnNames []string) string {
	builder := NewSelectQueryBuilder(tableName)
	builder.SetColumns(ConvertToColumns(columnNames))
	idColumn, _ := NewColumn("id", "")
	builder.SetWhere(tableName, *idColumn, "=", NewTokenValue("?"), "")
	return builder.Build()
}

func ReadAllQuery(tableName string, columnNames []string) string {
	builder := NewSelectQueryBuilder(tableName)
	builder.SetColumns(ConvertToColumns(columnNames))
	return builder.Build()
}

func UpdateQuery(tableName string, columnNames []string) string {
	// skip id
	builder := NewUpdateQueryBuilder(tableName)
	builder.SetNewValues(ConvertToColumns(columnNames), GetTokens(len(columnNames), "?"))
	idCol, _ := NewColumn("id", "")
	builder.SetWhere(tableName, *idCol, "=", NewTokenValue("?"), "")
	return builder.Build()
}

func DeleteQuery(tableName string) string {
	builder := NewDeleteQueryBuilder(tableName)
	idCol, _ := NewColumn("id", "")
	builder.SetWhere(tableName, *idCol, "=", NewTokenValue("?"), "")
	return builder.Build()
}
