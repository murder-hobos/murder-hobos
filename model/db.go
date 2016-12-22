package model

import (
	"database/sql"
	"errors"
	"fmt"

	//Import for side effects
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	// ErrInvalidID is raised when a given database ID is
	// either negative or does not exist
	ErrInvalidID = errors.New("model: invalid userID")
	// ErrNoResult is raised when a query returns no results
	ErrNoResult = sql.ErrNoRows
	// ErrNoResult is here to wrap the sql error. In our queries,
	// we can simply return the sql error, and the caller can
	// check if the error matches ErrNoResult
	// This hides the implementation of our physical datastore -
	// could use nosql/key-value store if we so chose

)

// Datastore defines methods for accessing a datastore containing information
// for our murder-hobos application.
//
// If errors are encountered for any method, the result should be nil along
// with the error value.
type Datastore interface {
	SpellDatastore
	ClassDatastore
	CharacterDatastore
	UserDatastore
}

// DB is a wrapper struct for our database connection that we
// use to implement the Datastore interface
type DB struct {
	*sqlx.DB
}

// NewDB returns an initialized DB connected to the mysql database
// described by the given DataSourceName. If an error is encountered,
// nil is returned, along with the error. func NewDB(dsn string) (*DB, error) {
func NewDB(dsn string) (Datastore, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// transact wraps an anonymous function with various db queries
// in a transaction
func (db *DB) transact(txFunc func(*sqlx.Tx) error) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	return txFunc(tx)
}
