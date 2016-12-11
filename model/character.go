package model

import (
	"log"

	"database/sql"
)

// CharacterDatastore describes methods available on our database
// pertaining to Characters
type CharacterDatastore interface {
	GetAllCharacters(userID int) (*[]Character, error)
	GetCharacterByName(userID int, name string) (*Character, error)
	InsertCharacter(char *Character) (int, error)
}

// Character represents our database Character table
type Character struct {
	ID                   int           `db:"id"`
	Name                 string        `db:"name"`
	Race                 string        `db:"race"`
	SpellAbilityModifier sql.NullInt64 `db:"spell_ability_modifier"`
	ProficiencyBonus     sql.NullInt64 `db:"proficiency_bonus"`
	UserID               int           `db:"user_id"`
}

// AbilityStr returns the character's ability modifier as a string if
// it is not NULL, otherwise returns "None"
func (c *Character) AbilityStr() string {
	if c.SpellAbilityModifier.Valid {
		return string(c.SpellAbilityModifier.Int64)
	}
	return "None"
}

// ProficiencyStr returns the character's proficiency bonus as a string if
// it is not NULL, otherwise returns "None"
func (c *Character) ProficiencyStr() string {
	if c.ProficiencyBonus.Valid {
		return string(c.ProficiencyBonus.Int64)
	}
	return "None"
}

// GetAllCharacters gets a list of every character belonging to a
// specified user
func (db *DB) GetAllCharacters(userID int) (*[]Character, error) {
	c := &[]Character{}
	err := db.Select(c, `SELECT id, name, race, spell_ability_modifier, proficiency_bonus, user_id 
					 FROM `+"`Character`"+` WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// GetCharacterByName get a list of every character that a user can view
func (db *DB) GetCharacterByName(userID int, name string) (*Character, error) {
	// verify arguments before hitting the db
	if name == "" {
		return nil, ErrNoResult
	}

	c := &Character{}
	err := db.Get(c, `SELECT * FROM `+"`Character`"+` WHERE user_id=? AND name=?`, userID, name)
	if err != nil {
		log.Printf("GetCharacterByName: %s\n", err.Error())
		return nil, err
	}

	return c, nil
}

// InsertCharacter inserts a given Character into the database.
// Returns the id of the new character if successful.
func (db *DB) InsertCharacter(char *Character) (int, error) {
	res, err := db.Exec(`INSERT INTO `+"`Character` "+`(name, race, spell_ability_modifier, proficiency_bonus,
						 user_id) VALUES (?, ?, ?, ?, ?)`,
		char.Name, char.Race, char.SpellAbilityModifier, char.ProficiencyBonus,
		char.UserID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
