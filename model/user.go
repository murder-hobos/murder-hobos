package model

// UserDatastore describes methods available to our database
// pertaining to users
type UserDatastore interface {
	GetUserByUsername(name string) (*User, bool)
	GetUserByID(id int) (*User, bool)
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
