package tales

import (
	"database/sql"
	"errors"
	"fmt"
	"talenest/backend/internal/data"
	"talenest/backend/internal/utils"
	"time"
)

const tableName = "tales"
const READ_BY_PARENT_STATEMENT = "READ_BY_PARENT"

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
	dbConn     *data.DatabaseConnector
	statements map[string]*sql.Stmt
}

func NewRepository(dbConn *data.DatabaseConnector) (Repository, error) {
	repo := &taleRepository{
		dbConn: dbConn,
	}
	var err error

	// Prepare statements only once (better perf + safety)
	createStmt, err := dbConn.PrepareQuery(
		data.CreateQuery(tableName, getColumnNames()))
	if err != nil {
		return nil, err
	}
	repo.statements = make(map[string]*sql.Stmt)
	repo.statements[data.CREATE_STATEMENT] = createStmt

	readByIdStmt, err := dbConn.PrepareQuery(
		data.ReadByIdQuery(tableName, getColumnNames()))
	if err != nil {
		return nil, err
	}
	repo.statements[data.READ_STATEMENT] = readByIdStmt

	readByParentIdStmt, err := dbConn.PrepareQuery(readByParentIdQuery())
	if err != nil {
		return nil, err
	}
	repo.statements[READ_BY_PARENT_STATEMENT] = readByParentIdStmt

	readAllStmt, err := dbConn.PrepareQuery(
		data.ReadAllQuery(tableName, getColumnNames()))
	if err != nil {
		return nil, err
	}
	repo.statements[data.READ_ALL_STATEMENT] = readAllStmt

	updateStmt, err := dbConn.PrepareQuery(
		data.UpdateQuery(tableName, getColumnNames()))
	if err != nil {
		return nil, err
	}
	repo.statements[data.UPDATE_STATEMENT] = updateStmt

	deleteStmt, err := dbConn.PrepareQuery(
		data.DeleteQuery(tableName))
	if err != nil {
		return nil, err
	}

	repo.statements[data.DELETE_STATEMENT] = deleteStmt

	return repo, nil
}

func getColumnNames() []string {
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

func (repo taleRepository) Create(tale *Tale) (int, error) {
	result, err := repo.statements[data.CREATE_STATEMENT].Exec(
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

func (repo taleRepository) ReadById(id int) (*Tale, error) {
	rows, err := repo.statements[data.READ_STATEMENT].Query(id)
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
	builder := data.NewSelectQueryBuilder(tableName)
	builder.SetColumns(data.ConvertToColumns(getColumnNames()))
	parentIdColumn, _ := data.NewColumn("parent_id", "")
	builder.SetWhere(tableName, *parentIdColumn, "=", data.NewTokenValue("?"), "")
	return builder.Build()
}

func (repo taleRepository) ReadByParentId(parentId int) (*Tales, error) {
	rows, err := repo.statements[READ_BY_PARENT_STATEMENT].Query(parentId)
	defer rows.Close()
	if err != nil {
		return &Tales{}, err
	}
	return repo.readTales(rows)
}

func (repo taleRepository) ReadAll() (*Tales, error) {
	rows, err := repo.statements[data.READ_ALL_STATEMENT].Query()
	defer rows.Close()
	if err != nil {
		return &Tales{}, err
	}
	return repo.readTales(rows)
}

func (repo taleRepository) Update(tale Tale) error {
	var result sql.Result
	var err error
	if !tale.deleted.IsZero() {
		result, err = repo.statements[data.UPDATE_STATEMENT].Exec(
			tale.Id,
			tale.Name,
			tale.Summary,
			tale.ParentId,
			tale.Status.Id,
			utils.CleanTime(tale.created),
			utils.CleanTime(time.Now()),
			utils.CleanTime(tale.deleted),
			tale.Id,
		)
	} else {
		result, err = repo.statements[data.UPDATE_STATEMENT].Exec(
			tale.Id,
			tale.Name,
			tale.Summary,
			tale.ParentId,
			tale.Status.Id,
			utils.CleanTime(tale.created),
			utils.CleanTime(time.Now()),
			nil,
			tale.Id,
		)
	}
	if err != nil {
		return err
	}
	if nRows, err := result.RowsAffected(); nRows != 1 || err != nil {
		return errors.New("The rows affected are different than 1")
	}
	return nil
}

func (repo taleRepository) Delete(id int) error {
	// permanent delete
	_, err := repo.statements[data.DELETE_STATEMENT].Exec(id)
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
