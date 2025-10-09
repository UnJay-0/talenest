package data

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

type DatabaseConnector struct {
	db *sql.DB
}

func NewDatabaseConnector(driver, dbPath, migrationPath string) *DatabaseConnector {
	db, err := sql.Open(driver, dbPath)
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

	m, err := migrate.New(
		"file://"+migrationPath,
		driver+"://"+dbPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	// connected
	return &DatabaseConnector{
		db: db,
	}
}

func (dbConnector *DatabaseConnector) Query(query string) *sql.Rows {
	rows, err := dbConnector.db.Query(query)
	if err != nil {
		fmt.Println(err)
		return &sql.Rows{}
	}
	return rows
}

func (dbConnector *DatabaseConnector) InsertQuery(query string, args ...Value) (int64, error) {
	result, err := dbConnector.ExecuteQuery(query, args)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return result.LastInsertId()
}

func (dbConnector *DatabaseConnector) ExecuteQuery(query string, args ...[]Value) (sql.Result, error) {
	result, err := dbConnector.db.Exec(query, args)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, err
}

func (dbConnector *DatabaseConnector) PrepareQuery(query string) (*sql.Stmt, error) {
	statement, err := dbConnector.db.Prepare(query)
	if err != nil {
		// TODO: Handle error
		fmt.Printf("Error preparing query: %v\n", err)
	}
	return statement, nil
}
