package model

import (
	"bytes"
	"database/sql"
	"html/template"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/murder-hobos/murder-hobos/util"
)

// SpellDatastore describes valid methods we have on our database
// pertaining to Spells
type SpellDatastore interface {
	GetAllCannonSpells() (*[]Spell, error)
	GetCannonSpellByName(name string) (*Spell, error)
	SearchCannonSpells(name string) (*[]Spell, error)
	FilterCannonSpells(level, school string) (*[]Spell, error)

	GetAllUserSpells(userID int) (*[]Spell, error)
	GetUserSpellByName(userID int, name string) (*Spell, error)
	SearchUserSpells(userID int, name string) (*[]Spell, error)
	FilterUserSpells(userID int, level, school string) (*[]Spell, error)

	GetSpellByID(id int) (*Spell, error)
	GetSpellClasses(spellID int) (*[]Class, error)
}

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

// LevelStr provides the spell's level as a string, with "Cantrip" for level 0
func (s *Spell) LevelStr() string {
	if s.Level == "0" {
		return "Cantrip"
	}
	return s.Level
}

// GetAllCannonSpells returns a list of every cannon spell object
// in our database (PHB, EE, SCAG)
func (db *DB) GetAllCannonSpells() (*[]Spell, error) {
	spells := &[]Spell{}
	if err := db.Select(spells, `SELECT * FROM CannonSpells`); err != nil {
		log.Printf("model: GetAllCannonSpells: %s", err.Error())
		return nil, err
	}
	return spells, nil
}

// GetAllUserSpells gets a list of every spell that a
// specified user has created in our database
func (db *DB) GetAllUserSpells(userID int) (*[]Spell, error) {
	if userID <= 0 {
		return nil, ErrInvalidID
	}

	spells := &[]Spell{}
	err := db.Select(spells, `SELECT * FROM Spell WHERE source_id=?`, userID)
	if err != nil {
		return nil, err
	}
	return spells, nil
}

// SearchCannonSpells gets a list of cannon spells with names similar
// to `name`
func (db *DB) SearchCannonSpells(name string) (*[]Spell, error) {
	query := `SELECT * FROM CannonSpells
			  WHERE name LIKE CONCAT ('%', ?, '%')
			  ORDER BY name ASC`
	spells := &[]Spell{}
	if err := db.Select(spells, query, name); err != nil {
		log.Printf("Error executing query %s\n %s\n", query, err.Error())
		return nil, err
	}

	return spells, nil
}

// SearchUserSpells gets a list of a user's spells with names similar
// to `name`
func (db *DB) SearchUserSpells(userID int, name string) (*[]Spell, error) {
	// don't hit the db with bunk query
	if userID <= 0 {
		return nil, ErrInvalidID
	}
	if name == "" {
		return nil, ErrNoResult
	}

	spells := &[]Spell{}
	err := db.Select(spells, `SELECT * FROM Spell 
							  WHERE source_id=? 
							  AND name LIKE CONCAT ('%', ?, '%')
							  ORDER BY name ASC;`, userID, name)
	if err != nil {
		log.Printf("model: SearchUserSpellByName: %s\n", err.Error())
		return nil, err
	}

	return spells, nil
}

// GetCannonSpellByName returns a single cannon spell with matching name
func (db *DB) GetCannonSpellByName(name string) (*Spell, error) {
	if name == "" {
		return nil, ErrNoResult
	}

	s := &Spell{}
	err := db.Get(s, "SELECT * FROM CannonSpells WHERE name=?", name)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetUserSpellByName returns a single cannon spell with matching name
func (db *DB) GetUserSpellByName(userID int, name string) (*Spell, error) {
	if userID <= 0 {
		return nil, ErrInvalidID
	}
	if name == "" {
		return nil, ErrNoResult
	}

	s := &Spell{}
	err := db.Select(s, "SELECT * FROM Spell WHERE source_id=? AND name=?", userID, name)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// FilterCannonSpells returns a list of cannon spells matching
// the search critera. If an empty argument is passed to one of the
// filters, that argument is not considered for filtering.
func (db *DB) FilterCannonSpells(level, school string) (*[]Spell, error) {
	if level == "" && school == "" {
		return nil, ErrNoResult
	}

	eqs := sq.Eq{}
	if level != "" {
		eqs["level"] = level
	}
	if school != "" {
		eqs["school"] = school
	}

	query, args, err := sq.Select("*").From("CannonSpells").Where(eqs).ToSql()

	spells := &[]Spell{}
	err = db.Select(spells, query, args...)
	if err != nil {
		return nil, err
	}
	return spells, nil
}

// FilterUserSpells returns a list of user spells matching
// the search critera. If an empty argument is passed to one of the
// filters, that argument is not considered for filtering.
// NOTE: name is given as a search param, not matched exactly
func (db *DB) FilterUserSpells(userID int, level, school string) (*[]Spell, error) {
	if userID <= 0 {
		return nil, ErrInvalidID
	}
	if level == "" && school == "" {
		return nil, ErrNoResult
	}

	eqs := sq.Eq{}
	eqs["source_id"] = userID

	if level != "" {
		eqs["level"] = level
	}
	if school != "" {
		eqs["school"] = school
	}

	query, args, err := sq.Select("*").From("Spell").Where(eqs).ToSql()

	spells := &[]Spell{}
	err = db.Select(spells, query, args...)
	if err != nil {
		return nil, err
	}
	return spells, nil
}

// // GetSchoolSpells searches the database and returns a slice of
// // Spell within same School id
// func (db *DB) GetSchoolSpells(school string, userID int, isCannon bool) (*Spell, error) {
// 	// verify arguments before hitting the db
// 	if school == "" {
// 		return nil, ErrNoResult
// 	}
// 	if userID < 0 {
// 		return nil, ErrInvalidID
// 	}
// 	if userID == 0 && !isCannon {
// 		return nil, ErrNoResult
// 	}
//
// 	var ids []int
// 	if userID > 0 { // If given a specific user, only search that
// 		ids = append(ids, userID)
// 	} else { // at this point isCannon must be true
// 		ids = append(ids, cannonIDs...)
// 	}
//
// 	query, args, err := sqlx.In(`SELECT * FROM Spell
// 								WHERE school=? AND
// 								source_id in (?);`,
// 		school, ids)
// 	if err != nil {
// 		log.Printf("Error preparing sqlx.In statement: %s\n", err.Error())
// 		return nil, err
// 	}
// 	query = db.Rebind(query)
//
// 	s := &Spell{}
// 	if err := db.Get(s, query, args...); err != nil {
// 		log.Printf("Error executing query %s\n %s\n", query, err.Error())
// 		return nil, err
// 	}
// 	return s, nil
// }

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

// GetSpellByID returns a single spell with matching id
func (db *DB) GetSpellByID(id int) (*Spell, error) {
	if id <= 0 {
		return nil, ErrNoResult
	}
	s := &Spell{}
	if err := db.Get(s, "SELECT * FROM Spell WHERE id=?", id); err != nil {
		return nil, err
	}
	return s, nil
}
