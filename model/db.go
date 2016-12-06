package model

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// Cannon ids
const (
	phbID  = 1
	eeID   = 2
	scagID = 3
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

	cannonIDs = []int{phbID, eeID, scagID}
)

// Datastore defines methods for accessing a datastore containing information
// for our murder-hobos application.
//
// If errors are encountered for any method, the result should be nil along
// with the error value.
//
// The rational behind using parameters to differentiate user spells and
// cannon spells is that we can see real performance benefits from sending
// 1 query across the network to our db instead of sending multiple and joining
// the results on the webserver.
//
//
// GetAllSpells returns all spells matching the arguments conditions.
// if userID is 0, no user-specific spells are included. If otherwise
// specified, search is restricted to that user.
// includeCannon decides whether or not to include cannon spells.
// Passing 0 to userID and false to includedCannon should return an ErrNoResult
// (hopefully obvious).
//
// GetSpellByID returns a single spell with specified id, as well as any
// errors encountered.
//
// GetSpellByName gets a single spell with the matching name. If userID
// is not 0, search is restricted to that user's spells.
// If isCannon is true, search is restricted to cannon spells.
// Passing a 0 userID and false for isCannon returns an ErrNoResult (hopefully obvious)
//
// GetSpellClasses gets a list of the Classes that a spell with the
// given id is available to.
type Datastore interface {
	GetAllSpells(userID int, includeCannon bool) (*[]Spell, error)
	GetSpellByID(id int) (*Spell, error)
	GetSpellByName(name string, userID int, isCannon bool) (*Spell, error)
	SearchSpellsByName(userID int, name string) (*[]Spell, error)

	GetSpellClasses(spellID int) (*[]Class, error)

	GetAllClasses() (*[]Class, error)
	GetClassByName(name string) (*Class, error)
	GetClassSpells(classID int) (*[]Spell, error)
}

// DB is a wrapper struct for our database connection that we
// use to implement the Datastore interface
type DB struct {
	*sqlx.DB
}

// NewDB returns an initialized DB connected to the mysql database
// described by the given DataSourceName. If an error is encountered,
// nil is returned, along with the error. func NewDB(dsn string) (*DB, error) {
func NewDB(dsn string) (*DB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
