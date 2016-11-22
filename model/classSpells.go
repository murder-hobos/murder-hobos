package model

// ClassSpells represents our db's ClassSpells table.
type ClassSpells struct {
	ClassID int `db:"class_id"`
	SpellID int `db:"spell_id"`
}
