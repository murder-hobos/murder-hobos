package model

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

// Class represents our database Class table
type Class struct {
	ID        int           `db:"id"`
	Name      string        `db:"name"`
	BaseClass sql.NullInt64 `db:"base_class_id"`
}

/*Gets all the classes
 */
func (db *DB) GetAllClasses(userID int) (*[]Class, error) {
	// verify arguments
	if userID == 0 {
		return nil, ErrNoResult
	}
	if userID < 0 {
		return nil, ErrInvalidID
	}

	var ids []int
	if userID > 0 {
		ids = append(ids, userID)
	}

	query, args, err := sqlx.In(`SELECT * FROM Class WHERE source_id IN (?);`, ids)
	if err != nil {
		log.Printf("Error preparing sqlx.In statement: %s\n", err.Error())
		return nil, err
	}
	query = db.Rebind(query)

	classes := &[]Class{}
	if err := db.Select(classes, query, args...); err != nil {
		log.Printf("Error executing query %s\n %s\n", query, err.Error())
		return nil, err
	}
	return classes, nil
}
