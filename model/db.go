package model

import (
	"database/sql"
	"errors"
	"log"

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
// is not 0, search should be restricted to that user's spells.
// If isCannon is true, search should be restricted to cannon spells.
// Passing a 0 userID and false for isCannon should return an
// ErrNoResult (hopefully obvious)
//
// GetSpellClasses gets a list of the Classes that a spell with the
// given id is available to.
type Datastore interface {
	GetAllSpells(userID int, includeCannon bool) (*[]Spell, error)
	GetSpellByID(id int) (*Spell, error)
	GetSpellByName(name string, userID int, isCannon bool) (*Spell, error)
	GetSpellClasses(spellID int) (*[]Class, error)
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

// GetAllSpells returns a slice of all spells in the database.
// userID can be 0. If otherwise specified, spells with the corresponding sourceID are
// included in the result.
// includeCannon chooses whether or not to include cannon(PHB, EE, SCAG) spells.
func (db *DB) GetAllSpells(userID int, includeCannon bool) (*[]Spell, error) {
	// verify arguments
	if userID == 0 && !includeCannon {
		return nil, ErrNoResult
	}
	if userID < 0 {
		return nil, ErrInvalidID
	}

	var ids []int
	if userID > 0 {
		ids = append(ids, userID)
	}
	if includeCannon {
		ids = append(ids, cannonIDs...)
	}

	query, args, err := sqlx.In(`SELECT * FROM Spell WHERE source_id IN (?);`, ids)
	if err != nil {
		log.Printf("Error preparing sqlx.In statement: %s\n", err.Error())
		return nil, err
	}
	query = db.Rebind(query)

	spells := &[]Spell{}
	if err := db.Select(spells, query, args...); err != nil {
		log.Printf("Error executing query %s\n %s\n", query, err.Error())
		return nil, err
	}
	return spells, nil
}

// GetSpellByID searches db for a Spell row with a matching id
func (db *DB) GetSpellByID(id int) (*Spell, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}
	s := &Spell{}
	if err := db.Get(s, "SELECT * FROM Spell WHERE id=?", id); err != nil {
		return nil, err
	}
	return s, nil
}

// GetSpellByName searches the datastore for a spell with the matching name.
// userID may be 0 for no user. If otherwise specified, search is restricted to that user's spells.
// includeCannon decides whether or not to search cannon spells.
// These options exist to enable different users to create spells with the same name,
// or the same name as cannon spells if they so choose
// NOTE: specifying nil userID and false for isCannon returns no result (hopefully obvious)
func (db *DB) GetSpellByName(name string, userID int, isCannon bool) (*Spell, error) {
	// verify arguments before hitting the db
	if name == "" {
		return nil, ErrNoResult
	}
	if userID < 0 {
		return nil, ErrInvalidID
	}
	if userID == 0 && !isCannon {
		return nil, ErrNoResult
	}

	var ids []int
	if userID > 0 { // If given a specific user, only search that
		ids = append(ids, userID)
	} else { // at this point isCannon must be true
		ids = append(ids, cannonIDs...)
	}

	query, args, err := sqlx.In(`SELECT * FROM Spell
								WHERE name=? AND
								source_id in (?);`,
		name, ids)
	if err != nil {
		log.Printf("Error preparing sqlx.In statement: %s\n", err.Error())
		return nil, err
	}
	query = db.Rebind(query)

	s := &Spell{}
	if err := db.Get(s, query, args...); err != nil {
		log.Printf("Error executing query %s\n %s\n", query, err.Error())
		return nil, err
	}
	return s, nil
}

// GetSpellClasses searches the database and returns a slice of
// Class objects available to the spell with spellID
func (db *DB) GetSpellClasses(spellID int) (*[]Class, error) {
	if spellID <= 0 {
		return nil, ErrNoResult
	}

	cs := &[]Class{}
	err := db.Select(cs, `SELECT C.id, C.name, C.base_class_id
	 					  FROM Class AS C
						  JOIN ClassSpells as CS ON
						  C.id = CS.class_id
						  JOIN Spell AS S ON
						  CS.spell_id = S.id
						  WHERE S.id = ?`,
		spellID)
	if err != nil {
		return nil, err
	}
	return cs, nil
}
