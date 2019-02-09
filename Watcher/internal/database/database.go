package database

import (
	"database/sql"
	"fmt"

	// Import PQ for the sql package
	_ "github.com/lib/pq"
)

// DB is a global entrypoint for database queries
var DB *Database

// Database interface wrapper for sql/pq
type Database struct {
	// Credentials
	Username string
	Password string
	Name     string

	// Connection settings
	URL  string
	Port int

	// Internal settings
	connection   *sql.DB
	setupQueries []string
}

// NewDB is a constructor for the database
func NewDB(username, password, name, url string, port int) *Database {
	db := Database{
		username,
		password,
		name,
		url,
		port,
		nil,
		[]string{},
	}
	db.Connect()

	DB = &db
	return &db
}

// Connect opens a connection to the dabase
func (db *Database) Connect() {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", db.Username, db.Password, db.Name)
	var err error

	// Connect
	db.connection, err = sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("Failed to connect to the postgres database")
		return
	}

	// Ping the connection
	err = db.connection.Ping()
	if err != nil {
		fmt.Println("Failed to ping the database server")
		return
	}

	fmt.Println("Succesfully connected to PostgresDB")
}

// AddTable adds a table to the setup sequence
func (db *Database) AddTable(query string) {
	db.setupQueries = append(db.setupQueries, query)
}

// CreateTables creates multiple tables based in the array of strings provided
func (db *Database) CreateTables() error {
	tx, err := db.connection.Begin()
	if err != nil {
		fmt.Println("Failed to create a transaction, " + err.Error())
	}

	for _, query := range db.setupQueries {
		_, err := tx.Exec(query)
		if err != nil {
			fmt.Println("Database", err.Error())
			continue
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Database", err.Error())
	}

	return nil
}

// Exec is short for db.connection.Exec()
func (db *Database) Exec(query string, args ...interface{}) {
	_, err := db.connection.Exec(query, args...)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// ExecBatch executes a query for every value in a BatchQuery struct
// it all happens in a transaction
func (db *Database) ExecBatch(batch BatchQuery) {
	tx, err := db.connection.Begin()
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, vars := range batch.Values {
		_, err = tx.Exec(batch.Query, vars...)
		if err != nil {
			fmt.Println(err.Error() + "\n" + batch.Query)
		}
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
	}
}

// Query is short for db.connection.query
func (db *Database) Query(sql string, vals ...interface{}) (*sql.Rows, error) {
	return db.connection.Query(sql, vals...)
}

// BatchQuery is a struct representing a query and a list of interfaces
// as context values
type BatchQuery struct {
	Query  string
	Values [][]interface{}
}

// AddValues adds a undetermined list of interfaces to the values
func (bq *BatchQuery) AddValues(vals ...interface{}) {
	bq.Values = append(bq.Values, vals)
}
