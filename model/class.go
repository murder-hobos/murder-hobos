package model

import (
	"log"

	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Class represents our database Class table
type Class struct {
	ID        int           `db:"id"`
	Name      string        `db:"name"`
	BaseClass sql.NullInt64 `db:"base_class_id"`
}

// GetAllClasses gets a list of every class in our database
func (db *DB) GetAllClasses(mainClass int) (*[]Class, error) {

	if mainClass != "" {
		query, args, err := sqlx.In(`SELECT * FROM Class 
									WHERE id = ?
									OR base_class_id = ?;`, mainClass, mainClass)
		if err != nil {
			log.Printf("Error preparing sqlx.In statement: %s\n", err.Error())
			return nil, err
		}
		query = db.Rebind(query)

		cs := &[]Class{}
		if err := db.Select(cs, query, args...); err != nil {
			log.Printf("Error executing query %s\n %s\n", query, err.Error())
			return nil, err
		}
		return cs, nil
	} else {
		query, args, err := sqlx.In(`SELECT id, name, base_class_id FROM Class`)
		if err != nil {
			log.Printf("Error preparing sqlx.In statement: %s\n", err.Error())
			return nil, err
		}
		query = db.Rebind(query)

		cs := &[]Class{}
		if err := db.Select(cs, query, args...); err != nil {
			log.Printf("Error executing query %s\n %s\n", query, err.Error())
			return nil, err
		}
		return cs, nil
	}
}

// GetClassByName get a list of every spells that a class can use
func (db *DB) GetClassByName(name string) (*Class, error) {
	// verify arguments before hitting the db
	if name == "" {
		return nil, ErrNoResult
	}

	c := &Class{}
	err := db.Get(c, `SELECT * FROM Class WHERE name=?`, name)
	if err != nil {
		log.Printf("GetClassByName: %s\n", err.Error())
		return nil, err
	}

	return c, nil
}

// GetClassSpells searches the database and returns a slice of
// Spell objects available to the class with classID
func (db *DB) GetClassSpells(classID int) (*[]Spell, error) {
	if classID <= 0 {
		return nil, ErrNoResult
	}

	spells := &[]Spell{}
	err := db.Select(spells, `SELECT S.id, S.name
	 					  FROM Spell AS S
						  JOIN ClassSpells as CS ON
						  S.id = CS.spell_id
						  JOIN Class AS C ON
						  CS.class_id = C.id
						  WHERE C.id = ?`,
		classID)
	if err != nil {
		return nil, err
	}
	return spells, nil
}
