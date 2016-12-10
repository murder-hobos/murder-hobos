package model

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// UserDatastore describes methods available to our database
// pertaining to users
type UserDatastore interface {
	GetUserByUsername(name string) (*User, bool)
	GetUserByID(id int) (*User, bool)
	CreateUser(name, password string) (*User, bool)
}

// User represents a user in our application
type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password []byte `db:"password"`
}

// GetUserByUsername returns a user object with matching
// username
func (db *DB) GetUserByUsername(name string) (*User, bool) {
	u := &User{}
	err := db.Get(u, `SELECT id, username, password 
					  FROM User
					  WHERE username = ?`,
		name)
	if err != nil {
		return nil, false
	}
	return u, true
}

// GetUserByID returns a user object with matchin
// id if found
func (db *DB) GetUserByID(id int) (*User, bool) {
	u := &User{}
	err := db.Get(u, `SELECT id, username, password
					  FROM User
					  WHERE id=?`, id)
	if err != nil {
		return nil, false
	}
	return u, true
}

// CreateUser adds a user to our database
func (db *DB) CreateUser(name, password string) (*User, bool) {
	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("CreateUser: %s", err.Error())
		return nil, false
	}
	res, err := db.Exec(`INSERT INTO `+"`User`"+` (username, password) VALUES (?, ?)`,
		name, p)
	if err != nil {
		log.Printf("CreateUser: %s", err.Error())
		return nil, false
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("CreateUser: %s", err.Error())
		return nil, false
	}
	user := &User{
		ID:       int(id),
		Username: name,
	}
	return user, true
}
