package model

import (
	"log"

	"database/sql"
	"github.com/yosssi/ace"
)

// CharacterDatastore describes methods available on our database
// pertaining to Characters
type CharacterDatastore interface {
	GetAllCharacters(userID int) (*[]Character, error)
	GetCharacterByName(userID int, name string) (*Character, error)
	CreateCharacter(userID int, *Character) (int, error)
}

// Character represents our database Character table
type Character struct {
	ID                   int           `db:"id"`
	Name                 string        `db:"name"`
	Race                 string        `db:"race"`
	SpellAbilityModifier sql.NullInt64 `db:"spell_ability_modifier"`
	ProficienyBonus      sql.NullInt64 `db:"proficiency_bonus"`
	UserID               int           `db:"user_id"`
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

func (db *DB) CreateCharacter(userID int, *Character) (int, error) {
	res, err := db.Exec(`INSERT INTO `+"`Character`")
}