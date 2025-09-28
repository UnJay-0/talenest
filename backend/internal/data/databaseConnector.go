package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"
)

var once sync.Once

type DatabaseConnector struct {
	db *sql.DB
}

var database *DatabaseConnector

func newDatabaseConnector() *DatabaseConnector {
	db, err := sql.Open("sqlite", "pathtodb")
	if err != nil {
		fmt.Println(err)
		// Handle error
		return nil
	}
	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Println(pingErr)
		// Handle error
	}

	// connected
	return &DatabaseConnector{
		db: db,
	}
}

func GetInstance() *DatabaseConnector {
	if database != nil {
		return database
	}
	once.Do(func() {
		database = newDatabaseConnector()
	})
	return database
}

func (dbConnector *DatabaseConnector) Query(query string) *sql.Rows {
	rows, err := dbConnector.db.Query(query)
	if err != nil {
		fmt.Println(err)
		return &sql.Rows{}
	}
	return rows
}

func (dbConnector *DatabaseConnector) InsertQuery(query string, idQuery string) (int, error) {
	err := dbConnector.ExecuteQuery(query)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	var lastInsertedId int
	if idQuery != "" {
		err = dbConnector.db.QueryRow(idQuery).Scan(&lastInsertedId)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
	}
	return lastInsertedId, err
}

func (dbConnector *DatabaseConnector) ExecuteQuery(query string) error {
	_, err := dbConnector.db.Exec(query)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (dbConnector *DatabaseConnector) PrepareQuery(query string) (*sql.Stmt, error) {
	statement, err := dbConnector.db.Prepare(query)
	if err != nil {
		// TODO: Handle error
		fmt.Printf("Error preparing query: %v\n", err)
	}
	return statement, err
}
