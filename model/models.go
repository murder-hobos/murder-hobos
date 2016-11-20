package model

import (
	"bytes"
	"context"
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/jaden-young/murder-hobos/util"
	"github.com/jmoiron/sqlx"
)

// DB is a wrapper for our db connection that we can use to define
// queries on as methods
type DB struct {
	DB *sqlx.DB
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

// Class represents our database Class table
type Class struct {
	ID        int           `db:"id"`
	Name      string        `db:"name"`
	BaseClass sql.NullInt64 `db:"base_class_id"`
}

// ClassSpells represents our db's ClassSpells table.
type ClassSpells struct {
	ClassID int `db:"class_id"`
	SpellID int `db:"spell_id"`
}

// User represents a user in our application
type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password []byte `db:"password"`
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

// WithDB wraps a http.HandlerFunc with access to this database
// connection by placing a reference with key "db" in the request's
// context
func (db *DB) WithDB(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "db", db)
		r = r.WithContext(ctx)
		fn(w, r)
	}
}

// GetSpellByID searches `db` for a Spell row
// with a matching id. If unsuccessful, an empty
// spell is returned along with false for ok.
func (db *DB) GetSpellByID(id int) (*Spell, bool) {
	s := &Spell{}
	if err := db.DB.Get(s, "SELECT * FROM Spell WHERE id=?", id); err != nil {
		return &Spell{}, false
	}
	return s, true
}

// GetSpellByName searches `db` for a Spell row
// with a matching `name` and a source_id in
// `sourceIDs`. If unsuccessful, an empty
// spell is returned along with false for ok.
func (db *DB) GetSpellByName(name string, sourceIDs []string) (*Spell, bool) {
	query, args, err := sqlx.In(`SELECT * FROM Spell
								WHERE name=? AND
								source_id in (?);`,
		name, sourceIDs)
	if err != nil {
		log.Printf("Error preparing sqlx.In statement: %s\n", err.Error())
		return &Spell{}, false
	}
	query = db.DB.Rebind(query)

	s := &Spell{}
	if err := db.DB.Get(s, query, args...); err != nil {
		log.Printf("Error executing query %s\n %s\n", query, err.Error())
		return &Spell{}, false
	}
	return s, true
}

// GetSpellClasses searches `db` and returns a slice of
// Class objects available to the spell with `spellID`
func (db *DB) GetSpellClasses(spellID int) (*[]Class, error) {
	cs := &[]Class{}
	err := db.DB.Select(cs, `SELECT C.id, C.name, C.base_class_id
	 					  FROM Class AS C
						  JOIN ClassSpells as CS ON
						  C.id = CS.class_id
						  JOIN Spell AS S ON
						  CS.spell_id = S.id
						  WHERE S.id = ?`,
		spellID)
	if err != nil {
		return &[]Class{}, err
	}
	return cs, nil
}

// GetAllSpells returns a slice of every spell object in the database
// with a source_id in `sourceIDs`
func (db *DB) GetAllSpells(sourceIDs []string) (*[]Spell, bool) {
	query, args, err := sqlx.In(`SELECT * FROM Spell WHERE source_id IN (?);`, sourceIDs)
	if err != nil {
		log.Printf("Error preparing sqlx.In statement: %s\n", err.Error())
		return &[]Spell{}, false
	}
	query = db.DB.Rebind(query)

	spells := &[]Spell{}
	if err := db.DB.Select(spells, query, args...); err != nil {
		log.Printf("Error executing query %s\n %s\n", query, err.Error())
		return &[]Spell{}, false
	}
	return spells, true
}
