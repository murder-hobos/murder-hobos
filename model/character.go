package model

import (
	"log"

	"database/sql"
)

// Class represents our database Class table
type Character struct {
	ID        			 sql.NullInt64  	`db:"id"`
	Name      			 string        		`db:"name"`
	Race 	  			 string	       		`db:"race"`
	SpellAbilityModifier int  		   		`db:"spell_ability_modifier"`
	ProficienyBonus		 int 		   		`db:"proficiency_bonus"`
	UserID				 sql.NullInt64 		`db:"user_id"`
}

// GetAllClasses gets a list of every class in our database
func (db *DB) GetAllCharacters() (*[]Character, error) {

	c := &[]Character{}
	if err := db.Select(c, `SELECT id, name, race, spell_ability_modifier, proficiency_bonus, user_id FROM Character;`); err != nil {
		return nil, err
	}
	return c, nil
}

// GetCharacterByName get a list of every character that a user can view
func (db *DB) GetCharacterByName(name string) (*Character, error) {
	// verify arguments before hitting the db
	if name == "" {
		return nil, ErrNoResult
	}

	c := &Character{}
	err := db.Get(c, `SELECT * FROM Character WHERE name=?`, name)
	if err != nil {
		log.Printf("GetCharacterByName: %s\n", err.Error())
		return nil, err
	}

	return c, nil
}