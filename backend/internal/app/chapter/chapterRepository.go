package chapter

import (
	"database/sql"
	"errors"
	"fmt"
	"talenest/backend/internal/data"
)

const tableName = "chapters"
const READ_BY_TALE_STATEMENT = "READ_BY_TALE"

type Repository interface {
	Create(chapter *Chapter) (int, error)
	ReadById(id int) (*Chapter, error)
	ReadByTale(tale int) (*Chapters, error)
	ReadAll() (*Chapters, error)
	Update(chapter Chapter) error
	Delete(id int) error
	Close() error
}

type chapterRepository struct {
	dbConn     *data.DatabaseConnector
	statements map[string]*sql.Stmt
}

func NewRepository(dbConn *data.DatabaseConnector) (Repository, error) {
	repo := &chapterRepository{
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

	readByParentIdStmt, err := dbConn.PrepareQuery(readByTaleQuery())
	if err != nil {
		return nil, err
	}
	repo.statements[READ_BY_TALE_STATEMENT] = readByParentIdStmt

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
		"content",
		"sentiment",
		"tale_id",
	}
}

func (repo chapterRepository) Create(chapter *Chapter) (int, error) {
	result, err := repo.statements[data.CREATE_STATEMENT].Exec(
		nil,
		chapter.Content,
		chapter.sentiment,
		chapter.TaleId,
	)

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	chapter.Id = int(id)
	return chapter.Id, err
}

func (repo chapterRepository) ReadById(id int) (*Chapter, error) {
	rows, err := repo.statements[data.READ_STATEMENT].Query(id)
	defer rows.Close()
	if err != nil {
		return &Chapter{}, err
	}
	if rows.Next() {
		chapter := Chapter{}
		err := rows.Scan(
			&chapter.Id,
			&chapter.Content,
			&chapter.sentiment,
			&chapter.TaleId,
		)

		if err != nil {
			return &Chapter{}, err
		}
		return &chapter, nil
	}
	return &Chapter{}, errors.New(fmt.Sprintf("Error on retrieving the row: %v", err))
}

func (repo chapterRepository) readChapters(rows *sql.Rows) (*Chapters, error) {
	chapterCollection := &Chapters{}
	for rows.Next() {
		chapter := Chapter{}
		err := rows.Scan(
			&chapter.Id,
			&chapter.Content,
			&chapter.sentiment,
			&chapter.TaleId,
		)

		if err != nil {
			return &Chapters{}, err
		}

		chapterCollection.Add(&chapter)
	}
	return chapterCollection, nil
}

func readByTaleQuery() string {
	builder := data.NewSelectQueryBuilder(tableName)
	builder.SetColumns(data.ConvertToColumns(getColumnNames()))
	parentIdColumn, _ := data.NewColumn("tale_id", "")
	builder.SetWhere(tableName, *parentIdColumn, "=", data.NewTokenValue("?"), "")
	return builder.Build()
}

func (repo chapterRepository) ReadByTale(taleId int) (*Chapters, error) {
	rows, err := repo.statements[READ_BY_TALE_STATEMENT].Query(taleId)
	defer rows.Close()
	if err != nil {
		return &Chapters{}, err
	}
	return repo.readChapters(rows)
}

func (repo chapterRepository) ReadAll() (*Chapters, error) {
	rows, err := repo.statements[data.READ_ALL_STATEMENT].Query()
	defer rows.Close()
	if err != nil {
		return &Chapters{}, err
	}
	return repo.readChapters(rows)
}

func (repo chapterRepository) Update(chapter Chapter) error {
	result, err := repo.statements[data.UPDATE_STATEMENT].Exec(
		chapter.Id,
		chapter.Content,
		chapter.sentiment,
		chapter.TaleId,
		chapter.Id,
	)
	if err != nil {
		return err
	}
	if nRows, err := result.RowsAffected(); nRows != 1 || err != nil {
		return errors.New("The rows affected are different than 1")
	}
	return nil
}

func (repo chapterRepository) Delete(id int) error {
	// permanent delete
	_, err := repo.statements[data.DELETE_STATEMENT].Exec(id)
	return err
}

func (repo chapterRepository) Close() error {
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
