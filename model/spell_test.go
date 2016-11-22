package model

import (
	"database/sql"
	"html/template"
	"reflect"
	"testing"
)

func TestSpell_ComponentsStr(t *testing.T) {
	type fields struct {
		ID            int
		Name          string
		Level         string
		School        string
		CastTime      string
		Duration      string
		Range         string
		Verbal        bool
		Somatic       bool
		Material      bool
		MaterialDesc  sql.NullString
		Concentration bool
		Ritual        bool
		Description   string
		SourceID      int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"VSM(stuff)",
			fields{Verbal: true, Somatic: true, Material: true, MaterialDesc: sql.NullString{
				Valid:  true,
				String: "Stuff",
			},
			},
			"V, S, M (Stuff)",
		},
		{
			"V",
			fields{Verbal: true},
			"V",
		},
		{
			"VSM",
			fields{Verbal: true, Somatic: true, Material: true, MaterialDesc: sql.NullString{
				Valid:  false,
				String: "",
			},
			},
			"V, S, M",
		},
		{
			"VS",
			fields{Verbal: true, Somatic: true},
			"V, S",
		},
	}
	for _, tt := range tests {
		s := &Spell{
			ID:            tt.fields.ID,
			Name:          tt.fields.Name,
			Level:         tt.fields.Level,
			School:        tt.fields.School,
			CastTime:      tt.fields.CastTime,
			Duration:      tt.fields.Duration,
			Range:         tt.fields.Range,
			Verbal:        tt.fields.Verbal,
			Somatic:       tt.fields.Somatic,
			Material:      tt.fields.Material,
			MaterialDesc:  tt.fields.MaterialDesc,
			Concentration: tt.fields.Concentration,
			Ritual:        tt.fields.Ritual,
			Description:   tt.fields.Description,
			SourceID:      tt.fields.SourceID,
		}
		if got := s.ComponentsStr(); got != tt.want {
			t.Errorf("%q. Spell.ComponentsStr() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

// The only thing we need to test is that the type conversion was made successfully
func TestSpell_HTMLDescription(t *testing.T) {
	type fields struct {
		ID            int
		Name          string
		Level         string
		School        string
		CastTime      string
		Duration      string
		Range         string
		Verbal        bool
		Somatic       bool
		Material      bool
		MaterialDesc  sql.NullString
		Concentration bool
		Ritual        bool
		Description   string
		SourceID      int
	}
	tests := []struct {
		name   string
		fields fields
		want   template.HTML
	}{
		{
			"stuff",
			fields{Description: "stuff"},
			"stuff",
		},
	}
	for _, tt := range tests {
		s := &Spell{
			ID:            tt.fields.ID,
			Name:          tt.fields.Name,
			Level:         tt.fields.Level,
			School:        tt.fields.School,
			CastTime:      tt.fields.CastTime,
			Duration:      tt.fields.Duration,
			Range:         tt.fields.Range,
			Verbal:        tt.fields.Verbal,
			Somatic:       tt.fields.Somatic,
			Material:      tt.fields.Material,
			MaterialDesc:  tt.fields.MaterialDesc,
			Concentration: tt.fields.Concentration,
			Ritual:        tt.fields.Ritual,
			Description:   tt.fields.Description,
			SourceID:      tt.fields.SourceID,
		}
		if got := s.HTMLDescription(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Spell.HTMLDescription() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
