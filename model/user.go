package model

// User represents a user in our application
type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password []byte `db:"password"`
}
