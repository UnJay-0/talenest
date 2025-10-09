package tales

import (
	"database/sql"
	"errors"
	"fmt"
	"talenest/backend/internal/data"
	"talenest/backend/internal/utils"
	"time"
)

const taleTableName = "tales"

type Repository interface {
	Create(tale *Tale) (int, error)
	ReadById(id int) (*Tale, error)
	ReadByParentId(parentId int) (*Tales, error)
	ReadAll() (*Tales, error)
	Update(tale Tale) error
	Delete(id int) error
	Close() error
}

type taleRepository struct {
	dbConn             *data.DatabaseConnector
	statements         map[string]*sql.Stmt
	readByIdStmt       *sql.Stmt
	readByParentIdStmt *sql.Stmt
	readAllStmt        *sql.Stmt
	updateStmt         *sql.Stmt
	deleteStmt         *sql.Stmt
}

func NewRepository(dbConn *data.DatabaseConnector) (Repository, error) {
	repo := &taleRepository{
		dbConn: dbConn,
	}
	var err error

	// Prepare statements only once (better perf + safety)
	createStmt, err := dbConn.PrepareQuery(createQuery())
	if err != nil {
		return nil, err
	}
	repo.statements = make(map[string]*sql.Stmt)
	repo.statements["create"] = createStmt

	readByIdStmt, err := dbConn.PrepareQuery(readByIdQuery())
	if err != nil {
		return nil, err
	}
	repo.statements["readById"] = readByIdStmt

	readByParentIdStmt, err := dbConn.PrepareQuery(readByParentIdQuery())
	if err != nil {
		return nil, err
	}
	repo.statements["readByParentId"] = readByParentIdStmt

	readAllStmt, err := dbConn.PrepareQuery(readAllQuery())
	if err != nil {
		return nil, err
	}
	repo.statements["readAll"] = readAllStmt

	updateStmt, err := dbConn.PrepareQuery(updateQuery())
	if err != nil {
		return nil, err
	}
	repo.statements["update"] = updateStmt

	deleteStmt, err := dbConn.PrepareQuery(deleteQuery())
	if err != nil {
		return nil, err
	}

	repo.statements["delete"] = deleteStmt

	return repo, nil
}

func getTalesColumnNames() []string {
	return []string{
		"id",
		"name",
		"summary",
		"parent_id",
		"status_id",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

func getTalesColumns() []data.Column {
	columns := []data.Column{}
	for _, columnName := range getTalesColumnNames() {
		column, _ := data.NewColumn(columnName, "")
		columns = append(columns, *column)
	}
	return columns
}

func createQuery() string {
	builder := data.NewInsertQueryBuilder(taleTableName)
	builder.SetColumns(getTalesColumns())
	builder.SetValues(data.GetTokens(len(getTalesColumnNames()), "?"))
	queryStr, _ := builder.Build()
	return queryStr
}

func (repo taleRepository) Create(tale *Tale) (int, error) {
	result, err := repo.statements["create"].Exec(
		nil,
		tale.Name,
		tale.Summary,
		tale.ParentId,
		tale.Status.Id,
		utils.CleanTime(tale.created),
		utils.CleanTime(tale.updated),
		nil,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	tale.Id = int(id)
	return tale.Id, err
}

func readByIdQuery() string {
	builder := data.NewSelectQueryBuilder(taleTableName)
	builder.SetColumns(getTalesColumns())
	idColumn, _ := data.NewColumn("id", "")
	builder.SetWhere(taleTableName, *idColumn, "=", data.NewTokenValue("?"), "")
	return builder.Build()
}

func (repo taleRepository) ReadById(id int) (*Tale, error) {
	rows, err := repo.statements["readById"].Query(id)
	defer rows.Close()
	if err != nil {
		return &Tale{}, err
	}
	if rows.Next() {
		tale := Tale{}
		var statusId int
		var createdString, updatedString string
		var deletedString sql.NullString
		err := rows.Scan(
			&tale.Id,
			&tale.Name,
			&tale.Summary,
			&tale.ParentId,
			&statusId, // TODO: Get status from its repository
			&createdString,
			&updatedString,
			&deletedString,
		)

		if err != nil {
			return &Tale{}, err
		}
		tale.setCreated(createdString)
		tale.setUpdated(updatedString)
		if deletedString.Valid {
			tale.setDeleted(deletedString.String)
		}

		return &tale, nil
	}
	return &Tale{}, errors.New(fmt.Sprintf("Error on retrieving the row: %v", err))
}

func (repo taleRepository) readTales(rows *sql.Rows) (*Tales, error) {
	taleCollection := &Tales{}
	for rows.Next() {
		tale := Tale{}
		var statusId int
		var createdString, updatedString string
		var deletedString sql.NullString
		err := rows.Scan(
			&tale.Id,
			&tale.Name,
			&tale.Summary,
			&tale.ParentId,
			&statusId, // TODO: Get status from its repository
			&createdString,
			&updatedString,
			&deletedString,
		)

		if err != nil {
			return &Tales{}, err
		}
		tale.setCreated(createdString)
		tale.setUpdated(updatedString)
		if deletedString.Valid {
			tale.setDeleted(deletedString.String)
		}
		taleCollection.Add(&tale)
	}
	return taleCollection, nil
}

func readByParentIdQuery() string {
	builder := data.NewSelectQueryBuilder(taleTableName)
	builder.SetColumns(getTalesColumns())
	parentIdColumn, _ := data.NewColumn("parent_id", "")
	builder.SetWhere(taleTableName, *parentIdColumn, "=", data.NewTokenValue("?"), "")
	return builder.Build()
}

func (repo taleRepository) ReadByParentId(parentId int) (*Tales, error) {
	rows, err := repo.statements["readByParentId"].Query(parentId)
	defer rows.Close()
	if err != nil {
		return &Tales{}, err
	}
	return repo.readTales(rows)
}

func readAllQuery() string {
	builder := data.NewSelectQueryBuilder(taleTableName)
	builder.SetColumns(getTalesColumns())
	return builder.Build()
}

func (repo taleRepository) ReadAll() (*Tales, error) {
	rows, err := repo.statements["readAll"].Query()
	defer rows.Close()
	if err != nil {
		return &Tales{}, err
	}
	return repo.readTales(rows)
}

func updateQuery() string {
	// skip id
	builder := data.NewUpdateQueryBuilder(taleTableName)
	builder.SetNewValues(getTalesColumns(), data.GetTokens(len(getTalesColumnNames()), "?"))
	idCol, _ := data.NewColumn("id", "")
	builder.SetWhere(taleTableName, *idCol, "=", data.NewTokenValue("?"), "")
	return builder.Build()
}

func (repo taleRepository) Update(tale Tale) error {
	result, err := repo.statements["update"].Exec(
		tale.Id,
		tale.Name,
		tale.Summary,
		tale.ParentId,
		tale.Status.Id,
		tale.created,
		utils.CleanTime(time.Now()),
		tale.deleted,
		tale.Id,
	)
	if err != nil {
		return err
	}
	if nRows, err := result.RowsAffected(); nRows != 1 || err != nil {
		fmt.Println(nRows, err)
		return errors.New("The rows affected are different than 1")
	}
	return nil
}

func deleteQuery() string {
	builder := data.NewDeleteQueryBuilder(taleTableName)
	idCol, _ := data.NewColumn("id", "")
	builder.SetWhere(taleTableName, *idCol, "=", data.NewTokenValue("?"), "")
	return builder.Build()
}

func (repo taleRepository) Delete(id int) error {
	// permanent delete
	_, err := repo.statements["delete"].Exec(id)
	return err
}

func (repo taleRepository) Close() error {
	var errs error
	for _, statement := range repo.statements {
		if statement != nil {
			if currentErr := statement.Close(); currentErr != nil {
				errs = errors.Join(currentErr)
			}
		}
	}
	return errs
}
