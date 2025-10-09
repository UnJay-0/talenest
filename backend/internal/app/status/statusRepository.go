package status

import (
	"database/sql"
	"errors"
	"fmt"
	"talenest/backend/internal/data"
)

const tableName = "status"

type Repository interface {
	Create(s *Status) (int, error)
	ReadById(id int) (*Status, error)
	ReadAll() ([]Status, error)
	Update(s Status) error
	Delete(id int) error
	Close() error
}

type statusRepository struct {
	dbConn     *data.DatabaseConnector
	statements map[string]*sql.Stmt
}

func NewRepository(dbConn *data.DatabaseConnector) (Repository, error) {
	repo := &statusRepository{
		dbConn: dbConn,
	}

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
		"color",
	}
}

func (repo statusRepository) Create(status *Status) (int, error) {
	result, err := repo.statements[data.CREATE_STATEMENT].Exec(
		nil,
		status.Name,
		status.color,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	status.Id = int(id)
	return status.Id, err
}

func (repo statusRepository) ReadById(id int) (*Status, error) {
	rows, err := repo.statements[data.READ_STATEMENT].Query(id)
	defer rows.Close()
	if err != nil {
		return &Status{}, err
	}
	if rows.Next() {
		status := Status{}
		err := rows.Scan(
			&status.Id,
			&status.Name,
			&status.color,
		)
		if err != nil {
			return &Status{}, err
		}
		return &status, nil
	}
	return &Status{}, errors.New(fmt.Sprintf("Error on retrieving the row: %v", err))
}

func (repo statusRepository) ReadAll() ([]Status, error) {
	rows, err := repo.statements[data.READ_ALL_STATEMENT].Query()
	defer rows.Close()
	if err != nil {
		return []Status{}, err
	}
	var statusCollection []Status
	for rows.Next() {
		status := Status{}
		err := rows.Scan(
			&status.Id,
			&status.Name,
			&status.color,
		)
		if err != nil {
			return []Status{}, err
		}
		statusCollection = append(statusCollection, status)
	}
	return statusCollection, nil
}

func (repo statusRepository) Update(status Status) error {
	result, err := repo.statements[data.UPDATE_STATEMENT].Exec(
		status.Id,
		status.Name,
		status.color,
		status.Id,
	)
	if err != nil {
		return err
	}
	if nRows, err := result.RowsAffected(); nRows != 1 || err != nil {
		return errors.New("The rows affected are different than 1")
	}
	return nil
}

func (repo statusRepository) Delete(id int) error {
	_, err := repo.statements[data.DELETE_STATEMENT].Exec(id)
	return err
}

func (repo statusRepository) Close() error {
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
