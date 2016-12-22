package model

import (
	"log"
	"strconv"

	"database/sql"

	"github.com/jmoiron/sqlx"
)

// CharacterDatastore describes methods available on our database
// pertaining to Characters
type CharacterDatastore interface {
	GetAllCharacters(userID int) (*[]Character, error)
	GetCharacterByName(userID int, name string) (*Character, error)
	InsertCharacter(char *Character) (int64, error)
}

// Character represents our database Character table
type Character struct {
	ID                   int           `db:"id"`
	Name                 string        `db:"name"`
	Race                 string        `db:"race"`
	SpellAbilityModifier sql.NullInt64 `db:"spell_ability_modifier"`
	ProficiencyBonus     sql.NullInt64 `db:"proficiency_bonus"`
	UserID               int           `db:"user_id"`
	Levels               []ClassLevelView
}

// ClassLevelView is a single entry listing a character's level in
// a single class
type ClassLevelView struct {
	Class
	Level int `db:"level"`
}

// AbilityStr returns the character's ability modifier as a string if
// it is not NULL, otherwise returns "None"
func (c *Character) AbilityStr() string {
	if c.SpellAbilityModifier.Valid {
		return strconv.FormatInt(c.SpellAbilityModifier.Int64, 10)
	}
	return "None"
}

// ProficiencyStr returns the character's proficiency bonus as a string if
// it is not NULL, otherwise returns "None"
func (c *Character) ProficiencyStr() string {
	if c.ProficiencyBonus.Valid {
		return strconv.FormatInt(c.ProficiencyBonus.Int64, 10)
	}
	return "None"
}

// GetAllCharacters gets a list of every character belonging to a
// specified user
//
// FIXME: Right now GetAllCharacters does not populate the character's
// classlevels
func (db *DB) GetAllCharacters(userID int) (*[]Character, error) {
	c := &[]Character{}
	err := db.Select(c, `SELECT * FROM `+"`Character`"+` AS C
						 WHERE C.user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// InsertCharacter inserts a given Character into the database.
// Insertions are done inside of a transaction.
// Returns the id of the new character if successful.
func (db *DB) InsertCharacter(char *Character) (id int64, err error) {
	id = 0
	tx, err := db.Begin()
	if err != nil {
		return
	}

	db.transact(func(*sqlx.Tx) error {
		// Insert actual character
		res, err := tx.Exec(`INSERT INTO `+"`Character` "+`(name, race, spell_ability_modifier, proficiency_bonus,
						 user_id) VALUES (?, ?, ?, ?, ?)`,
			char.Name, char.Race, char.SpellAbilityModifier, char.ProficiencyBonus,
			char.UserID)
		if err != nil {
			return err
		}

		// set last insert id for our parent function's return value
		id, err = res.LastInsertId()
		if err != nil {
			return err
		}

		// Insert character levels
		stmt, err := tx.Prepare(`INSERT INTO CharacterLevels (char_id, class_id, level) VALUES (?, ?, ?)`)
		if err != nil {
			return err
		}

		for _, c := range char.Levels {
			_, err = stmt.Exec(id, c.Class.ID, c.Level)
			if err != nil {
				return err
			}
		}
		return err
	})
	return id, err
}

// GetCharacterByID returns a struct containing a nested Character
// as well as
func (db *DB) GetCharacterByID(charID int) (*Character, error) {
	var c *Character

	db.transact(func(tx *sqlx.Tx) error {
		err := tx.Get(c, `SELECT * FROM `+"`Character`"+` WHERE id=?`, charID)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		var clv *[]ClassLevelView
		err = tx.Select(clv, `SELECT id, name, level FROM LevelsForCharactersView WHERE char_id=?`, charID)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		c.Levels = *clv
		return nil
	})

	return c, nil
}

// GetCharacterByName get a list of every character that a user can view
func (db *DB) GetCharacterByName(userID int, name string) (*Character, error) {
	// verify arguments before hitting the db
	if name == "" {
		return nil, ErrNoResult
	}

	c := &Character{}
	db.transact(func(tx *sqlx.Tx) error {
		err := tx.Get(c, `SELECT * FROM `+"`Character`"+` 
						  WHERE user_id=? AND name=?`, userID, name)
		if err != nil {
			log.Printf("GetCharacterByName: %s\n", err.Error())
			return err
		}

		clv := []ClassLevelView{}
		err = tx.Select(&clv, `SELECT L.id, L.name, L.level, L.base_class_id
							FROM LevelsForCharactersView AS L
							JOIN `+"`Character` AS C"+`
							ON L.char_id = C.id
							WHERE C.user_id=? AND C.name=?`, userID, name)
		if err != nil {
			log.Println("GetCharacterByName")
			log.Println(err.Error())
			return err
		}
		c.Levels = clv
		return nil
	})
	return c, nil
}
