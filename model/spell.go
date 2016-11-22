package model

import (
	"bytes"
	"database/sql"
	"html/template"

	"github.com/jaden-young/murder-hobos/util"
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
