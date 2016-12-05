package model

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
