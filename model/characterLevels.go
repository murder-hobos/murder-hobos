package model

// CharacterLevels represents our database table we use
// for keeping track of character's levels in different
// classes
type CharacterLevels struct {
	CharID  int `db:"char_id"`
	ClassID int `db:"class_id"`
	Level   int `db:"level"`
}
