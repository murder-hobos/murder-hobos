package model

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestCharacter_AbilityStr(t *testing.T) {
	type fields struct {
		ID                   int
		Name                 string
		Race                 string
		SpellAbilityModifier sql.NullInt64
		ProficiencyBonus     sql.NullInt64
		UserID               int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"Valid",
			fields{
				SpellAbilityModifier: sql.NullInt64{
					Int64: 42,
					Valid: true,
				},
			},
			"42",
		},
		{
			"Invalid",
			fields{
				SpellAbilityModifier: sql.NullInt64{
					Int64: 0,
					Valid: false,
				},
			},
			"None",
		},
	}
	for _, tt := range tests {
		c := &Character{
			ID:                   tt.fields.ID,
			Name:                 tt.fields.Name,
			Race:                 tt.fields.Race,
			SpellAbilityModifier: tt.fields.SpellAbilityModifier,
			ProficiencyBonus:     tt.fields.ProficiencyBonus,
			UserID:               tt.fields.UserID,
		}
		if got := c.AbilityStr(); got != tt.want {
			t.Errorf("%q. Character.AbilityStr() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCharacter_ProficiencyStr(t *testing.T) {
	type fields struct {
		ID                   int
		Name                 string
		Race                 string
		SpellAbilityModifier sql.NullInt64
		ProficiencyBonus     sql.NullInt64
		UserID               int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"Valid",
			fields{
				ProficiencyBonus: sql.NullInt64{
					Int64: 42,
					Valid: true,
				},
			},
			"42",
		},
		{
			"Invalid",
			fields{
				ProficiencyBonus: sql.NullInt64{
					Int64: 0,
					Valid: false,
				},
			},
			"None",
		},
	}
	for _, tt := range tests {
		c := &Character{
			ID:                   tt.fields.ID,
			Name:                 tt.fields.Name,
			Race:                 tt.fields.Race,
			SpellAbilityModifier: tt.fields.SpellAbilityModifier,
			ProficiencyBonus:     tt.fields.ProficiencyBonus,
			UserID:               tt.fields.UserID,
		}
		if got := c.ProficiencyStr(); got != tt.want {
			t.Errorf("%q. Character.ProficiencyStr() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestDB_GetAllCharacters(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		userID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]Character
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		db := &DB{
			DB: tt.fields.DB,
		}
		got, err := db.GetAllCharacters(tt.args.userID)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. DB.GetAllCharacters() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. DB.GetAllCharacters() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestDB_GetCharacterByName(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		userID int
		name   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Character
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		db := &DB{
			DB: tt.fields.DB,
		}
		got, err := db.GetCharacterByName(tt.args.userID, tt.args.name)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. DB.GetCharacterByName() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. DB.GetCharacterByName() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestDB_InsertCharacter(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		char *Character
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		db := &DB{
			DB: tt.fields.DB,
		}
		got, err := db.InsertCharacter(tt.args.char)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. DB.InsertCharacter() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. DB.InsertCharacter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
