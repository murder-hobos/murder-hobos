package model

import "database/sql"

// Class represents our database Class table
type Class struct {
	ID        int           `db:"id"`
	Name      string        `db:"name"`
	BaseClass sql.NullInt64 `db:"base_class_id"`
}

// GetAllClasses gets a list of every class in our database
func (db *DB) GetAllClasses() (*[]Class, error) {

	cs := &[]Class{}
	if err := db.Select(cs, `SELECT * FROM Class`); err != nil {
		return nil, err
	}
	return cs, nil
}
