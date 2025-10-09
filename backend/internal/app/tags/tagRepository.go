package tags

import (
	"database/sql"
	"errors"
	"fmt"
	"talenest/backend/internal/data"
)

const tableName = "tag"

type Repository interface {
	Create(tag *Tag) (int, error)
	ReadById(id int) (*Tag, error)
	ReadAll() ([]Tag, error)
	Update(tag Tag) error
	Delete(id int) error
	Close() error
}

type tagRepository struct {
	dbConn     *data.DatabaseConnector
	statements map[string]*sql.Stmt
}

func getColumnNames() []string {
	return []string{
		"id",
		"name",
	}
}

func NewRepository(dbConn *data.DatabaseConnector) (Repository, error) {
	repo := &tagRepository{
		dbConn: dbConn,
	}

	var err error

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

func (repo tagRepository) Create(tag *Tag) (int, error) {
	result, err := repo.statements[data.CREATE_STATEMENT].Exec(
		nil,
		tag.Name,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	tag.Id = int(id)
	return tag.Id, err
}

func (repo tagRepository) ReadById(id int) (*Tag, error) {
	rows, err := repo.statements[data.READ_ALL_STATEMENT].Query(id)
	defer rows.Close()
	if err != nil {
		return &Tag{}, err
	}
	if rows.Next() {
		tag := Tag{}
		err := rows.Scan(
			&tag.Id,
			&tag.Name,
		)
		if err != nil {
			return &Tag{}, err
		}
		return &tag, nil
	}
	return &Tag{}, errors.New(fmt.Sprintf("Error on retrieving the row: %v", err))
}

func (repo tagRepository) ReadAll() ([]Tag, error) {
	rows, err := repo.statements[data.READ_ALL_STATEMENT].Query()
	defer rows.Close()
	if err != nil {
		return []Tag{}, err
	}
	tags := []Tag{}
	for rows.Next() {
		tag := Tag{}
		err := rows.Scan(
			&tag.Id,
			&tag.Name,
		)
		if err != nil {
			return []Tag{}, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (repo tagRepository) Update(tag Tag) error {
	result, err := repo.statements[data.UPDATE_STATEMENT].Exec(
		tag.Id,
		tag.Name,
		tag.Id,
	)
	if err != nil {
		return err
	}
	if nRows, err := result.RowsAffected(); nRows != 1 || err != nil {
		return errors.New("The rows affected are different than 1")
	}
	return nil
}

func (repo tagRepository) Delete(id int) error {
	_, err := repo.statements[data.DELETE_STATEMENT].Exec(id)
	return err
}

func (repo tagRepository) Close() error {
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
