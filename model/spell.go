package model

import (
	"bytes"
	"database/sql"
	"html/template"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/murder-hobos/murder-hobos/util"
)

// Spell represents our database version of a spell
type Spell struct {
	ID            int            `db:"id"`
	Name          string         `db:"name"`
	Level         string         `db:"level"`
	School        string         `db:"school"`
	CastTime      string         `db:"cast_time"`
	Duration      string         `db:"duration"`
	Range         string         `db:"range"`
	Verbal        bool           `db:"comp_verbal"`
	Somatic       bool           `db:"comp_somatic"`
	Material      bool           `db:"comp_material"`
	MaterialDesc  sql.NullString `db:"material_desc"`
	Concentration bool           `db:"concentration"`
	Ritual        bool           `db:"ritual"`
	Description   string         `db:"description"`
	SourceID      int            `db:"source_id"`
}

// ComponentsStr returns a string representation of the
// components for a spell.
// Example:
// 		"V, S, M (Some cool component no one will ever need because they have a focus)"
func (s *Spell) ComponentsStr() string {
	b := bytes.Buffer{}
	if s.Verbal {
		b.WriteString("V")
	}
	if s.Somatic {
		b.WriteString("S")
	}
	if s.Material {
		b.WriteString("M")
	}

	str := util.Intersperse(b.String(), ", ")
	b.Reset()

	// if not null
	if s.MaterialDesc.Valid {
		b.WriteString(" (")
		b.WriteString(s.MaterialDesc.String)
		b.WriteString(")")
	}
	return str + b.String()
}

// HTMLDescription returns the spell's Description as a template.HTML
// so that the HTML in the description will be rendered instead of
// escaped in go's templating engine
func (s *Spell) HTMLDescription() template.HTML {
	return template.HTML(s.Description)
}

// GetAllSpells returns a slice of all spells in the database.
// userID can be 0. If otherwise specified, spells with the corresponding sourceID are
// included in the result.
// includeCannon chooses whether or not to include cannon(PHB, EE, SCAG) spells.
func (db *DB) GetAllSpells(userID int, includeCannon bool) (*[]Spell, error) {
	// verify arguments
	if userID == 0 && !includeCannon {
		return nil, ErrNoResult
	}
	if userID < 0 {
		return nil, ErrInvalidID
	}

	var ids []int
	if userID > 0 {
		ids = append(ids, userID)
	}
	if includeCannon {
		ids = append(ids, cannonIDs...)
	}

	query, args, err := sqlx.In(`SELECT * FROM Spell WHERE source_id IN (?);`, ids)
	if err != nil {
		log.Printf("Error preparing sqlx.In statement: %s\n", err.Error())
		return nil, err
	}
	query = db.Rebind(query)

	spells := &[]Spell{}
	if err := db.Select(spells, query, args...); err != nil {
		log.Printf("Error executing query %s\n %s\n", query, err.Error())
		return nil, err
	}
	return spells, nil
}

// GetSpellByID searches db for a Spell row with a matching id
func (db *DB) GetSpellByID(id int) (*Spell, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}
	s := &Spell{}
	if err := db.Get(s, "SELECT * FROM Spell WHERE id=?", id); err != nil {
		return nil, err
	}
	return s, nil
}

// GetSpellByName searches the datastore for a spell with the matching name.
// userID may be 0 for no user. If otherwise specified, search is restricted to that user's spells.
// includeCannon decides whether or not to search cannon spells.
// These options exist to enable different users to create spells with the same name,
// or the same name as cannon spells if they so choose
// NOTE: specifying nil userID and false for isCannon returns no result (hopefully obvious)
func (db *DB) GetSpellByName(name string, userID int, isCannon bool) (*Spell, error) {
	// verify arguments before hitting the db
	if name == "" {
		return nil, ErrNoResult
	}
	if userID < 0 {
		return nil, ErrInvalidID
	}
	if userID == 0 && !isCannon {
		return nil, ErrNoResult
	}

	var ids []int
	if userID > 0 { // If given a specific user, only search that
		ids = append(ids, userID)
	} else { // at this point isCannon must be true
		ids = append(ids, cannonIDs...)
	}

	query, args, err := sqlx.In(`SELECT * FROM Spell
								WHERE name=? AND
								source_id IN (?);`,
		name, ids)
	if err != nil {
		log.Printf("Error preparing sqlx.In statement: %s\n", err.Error())
		return nil, err
	}
	query = db.Rebind(query)

	s := &Spell{}
	if err := db.Get(s, query, args...); err != nil {
		log.Printf("Error executing query %s\n %s\n", query, err.Error())
		return nil, err
	}
	return s, nil
}

// GetSpellClasses searches the database and returns a slice of
// Class objects available to the spell with spellID
func (db *DB) GetSpellClasses(spellID int) (*[]Class, error) {
	if spellID <= 0 {
		return nil, ErrNoResult
	}

	cs := &[]Class{}
	err := db.Select(cs, `SELECT C.id, C.name, C.base_class_id
	 					  FROM Class AS C
						  JOIN ClassSpells as CS ON
						  C.id = CS.class_id
						  JOIN Spell AS S ON
						  CS.spell_id = S.id
						  WHERE S.id = ?`,
		spellID)
	if err != nil {
		return nil, err
	}
	return cs, nil
}
