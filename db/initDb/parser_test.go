package initDb

import (
	"database/sql"
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
	"github.com/murder-hobos/murder-hobos/model"
)

func TestXMLSpell_ToDbSpell(t *testing.T) {
	type fields struct {
		XMLName    xml.Name
		Name       string
		Level      string
		School     string
		Ritual     string
		Time       string
		Range      string
		Components string
		Duration   string
		Classes    string
		Texts      []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    model.Spell
		wantErr bool
	}{
		{
			"Zone of truth. Should work.",
			fields{
				Name:       "Zone of Truth",
				Level:      "2",
				School:     "EN",
				Ritual:     "",
				Time:       "1 action",
				Range:      "60 feet",
				Components: "V, S, M",
				Duration:   "10 minutes",
				Classes:    "Bard, Cleric, Paladin, Paladin (Devotion), Paladin (Crown)",
				Texts: []string{"You create a magical zone that guards against deception in a 15-foot-radius sphere centered on a point of your choice within range. Until the spell ends, a creature that enters the spell's area for the first time on a turn or starts its turn there must make a Charisma saving throw. On a failed save, a creature can't speak a deliberate lie while in the radius. You know whether each creature succeeds or fails on its saving throw.",
					"",
					"An affected creature is aware of the spell and can thus avoid answering questions to which it would normally respond with a lie. Such creatures can be evasive in its answers as long as it remains within the boundaries of the truth.",
				},
			},
			model.Spell{
				Name:     "Zone of Truth",
				Level:    "2",
				School:   "Enchantment",
				CastTime: "1 action",
				Duration: "10 minutes",
				Range:    "60 feet",
				Verbal:   true,
				Somatic:  true,
				Material: true,
				MaterialDesc: sql.NullString{
					String: "",
					Valid:  false,
				},
				Concentration: false,
				Ritual:        false,
				Description:   "You create a magical zone that guards against deception in a 15-foot-radius sphere centered on a point of your choice within range. Until the spell ends, a creature that enters the spell&#39;s area for the first time on a turn or starts its turn there must make a Charisma saving throw. On a failed save, a creature can&#39;t speak a deliberate lie while in the radius. You know whether each creature succeeds or fails on its saving throw.<br/><br/>An affected creature is aware of the spell and can thus avoid answering questions to which it would normally respond with a lie. Such creatures can be evasive in its answers as long as it remains within the boundaries of the truth.",
				SourceID:      PHBid,
			},
			false,
		},
		{
			"Absorb Elements. Should be EE.",
			fields{
				Name:       "Absorb Elements (EE)",
				Level:      "1",
				School:     "A",
				Ritual:     "",
				Time:       "1 reaction, which you take when you take acid, cold, fire, lightning, or thunder damage",
				Range:      "",
				Components: "S",
				Duration:   "1 round",
				Classes:    "Druid, Ranger, Wizard, Fighter (Eldritch Knight)",
				Texts: []string{"The spell captures some of the incoming energy, lessening its effect on you and storing it for your next melee attack. You have resistance to the triggering damage type until the start of your next turn. Also, the first time you hit with a melee attack on your next turn, the target takes an extra 1d6 damage of the triggering type, and the spell ends.",
					"",
					"At Higher Levels: When you cast this spell using a spell slot of 2nd level or higher, the extra damage increases by 1d6 for each slot level above 1st.",
					"",
					"This spell can be found in the Elemental Evil Player's Companion",
				},
			},
			model.Spell{
				Name:     "Absorb Elements",
				Level:    "1",
				School:   "Abjuration",
				CastTime: "1 reaction, which you take when you take acid, cold, fire, lightning, or thunder damage",
				Duration: "1 round",
				Range:    "",
				Verbal:   false,
				Somatic:  true,
				Material: false,
				MaterialDesc: sql.NullString{
					String: "",
					Valid:  false,
				},
				Concentration: false,
				Ritual:        false,
				Description:   "The spell captures some of the incoming energy, lessening its effect on you and storing it for your next melee attack. You have resistance to the triggering damage type until the start of your next turn. Also, the first time you hit with a melee attack on your next turn, the target takes an extra 1d6 damage of the triggering type, and the spell ends.<br/><br/>At Higher Levels: When you cast this spell using a spell slot of 2nd level or higher, the extra damage increases by 1d6 for each slot level above 1st.<br/><br/>This spell can be found in the Elemental Evil Player&#39;s Companion",
				SourceID:      EEid,
			},
			false,
		},
	}
	for _, tt := range tests {
		x := XMLSpell{
			Name:       tt.fields.Name,
			Level:      tt.fields.Level,
			School:     tt.fields.School,
			Ritual:     tt.fields.Ritual,
			Time:       tt.fields.Time,
			Range:      tt.fields.Range,
			Components: tt.fields.Components,
			Duration:   tt.fields.Duration,
			Classes:    tt.fields.Classes,
			Texts:      tt.fields.Texts,
		}
		got, err := x.ToDbSpell()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. XMLSpell.ToDbSpell() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			spew.Dump(got, tt.want)
			t.Errorf("%q. XMLSpell.ToDbSpell() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestXMLSpell_ParseClasses(t *testing.T) {
	type fields struct {
		Name       string
		Level      string
		School     string
		Ritual     string
		Time       string
		Range      string
		Components string
		Duration   string
		Classes    string
		Texts      []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []model.Class
		want1  bool
	}{
		{
			"3 classes",
			fields{
				Classes: "Cleric, Cleric (Arcana), Druid",
			},
			[]model.Class{
				model.Class{
					ID:   2,
					Name: "Cleric",
					BaseClass: sql.NullInt64{
						Int64: 0,
						Valid: false,
					},
				},
				model.Class{
					ID:   3,
					Name: "Cleric (Arcana)",
					BaseClass: sql.NullInt64{
						Int64: 2,
						Valid: true,
					},
				},
				model.Class{
					ID:   12,
					Name: "Druid",
					BaseClass: sql.NullInt64{
						Int64: 0,
						Valid: false,
					},
				},
			},
			true,
		},
		{
			"No classes",
			fields{},
			[]model.Class{},
			false,
		},
		{
			"One class",
			fields{
				Classes: "Bard",
			},
			[]model.Class{
				model.Class{
					ID:   1,
					Name: "Bard",
					BaseClass: sql.NullInt64{
						Int64: 0,
						Valid: false,
					},
				},
			},
			true,
		},
		{
			"Absorb Elements (EE)",
			fields{
				Classes: "Druid, Ranger, Wizard, Fighter (Eldritch Knight)",
			},
			[]model.Class{
				model.Class{
					ID:   12,
					Name: "Druid",
					BaseClass: sql.NullInt64{
						Int64: 0,
						Valid: false,
					},
				},
				model.Class{
					ID:   27,
					Name: "Ranger",
					BaseClass: sql.NullInt64{
						Int64: 0,
						Valid: false,
					},
				},
				model.Class{
					ID:   34,
					Name: "Wizard",
					BaseClass: sql.NullInt64{
						Int64: 0,
						Valid: false,
					},
				},
				model.Class{
					ID:   36,
					Name: "Fighter (Eldritch Knight)",
					BaseClass: sql.NullInt64{
						Int64: 35,
						Valid: true,
					},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		x := &XMLSpell{
			Name:       tt.fields.Name,
			Level:      tt.fields.Level,
			School:     tt.fields.School,
			Ritual:     tt.fields.Ritual,
			Time:       tt.fields.Time,
			Range:      tt.fields.Range,
			Components: tt.fields.Components,
			Duration:   tt.fields.Duration,
			Classes:    tt.fields.Classes,
			Texts:      tt.fields.Texts,
		}
		got, got1 := x.ParseClasses()
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. XMLSpell.ParseClasses() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. XMLSpell.ParseClasses() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}

func TestXMLSpell_parseComponents(t *testing.T) {
	type fields struct {
		Name       string
		Level      string
		School     string
		Ritual     string
		Time       string
		Range      string
		Components string
		Duration   string
		Classes    string
		Texts      []string
	}
	tests := []struct {
		name   string
		fields fields
		want   components
	}{

		{
			"V",
			fields{Components: "V"},
			components{Verb: true, Som: false, Mat: false, Matdesc: sql.NullString{
				String: "",
				Valid:  false},
			},
		},
		{
			"V, S",
			fields{Components: "V, S"},
			components{Verb: true, Som: true, Mat: false, Matdesc: sql.NullString{
				String: "",
				Valid:  false},
			},
		},
		{
			"V, S, M",
			fields{Components: "V, S, M"},
			components{Verb: true, Som: true, Mat: true, Matdesc: sql.NullString{
				String: "",
				Valid:  false},
			},
		},
		{
			"V, S, M (text)",
			fields{Components: "V, S, M (a jade circlet worth at least 1,500 gp, which you must place on your head before you cast the spell)"},
			components{Verb: true, Som: true, Mat: true, Matdesc: sql.NullString{
				String: "A jade circlet worth at least 1,500 gp, which you must place on your head before you cast the spell",
				Valid:  true},
			},
		},

		{
			"S",
			fields{Components: "S"},
			components{Verb: false, Som: true, Mat: false, Matdesc: sql.NullString{
				String: "",
				Valid:  false},
			},
		},
		{
			"V, M (text)",
			fields{Components: "V, M (bat fur and a drop of pitch or piece of coal)"},
			components{Verb: true, Som: false, Mat: true, Matdesc: sql.NullString{
				String: "Bat fur and a drop of pitch or piece of coal",
				Valid:  true},
			},
		},
		{
			"S, M",
			fields{Components: "S, M"},
			components{Verb: false, Som: true, Mat: true, Matdesc: sql.NullString{
				String: "",
				Valid:  false},
			},
		},
		{
			"S, M (text)",
			fields{Components: "S, M (a glowing stick of incense or a crystal vial filled with phosphorescent material)"},
			components{Verb: false, Som: true, Mat: true, Matdesc: sql.NullString{
				String: "A glowing stick of incense or a crystal vial filled with phosphorescent material",
				Valid:  true},
			},
		},
	}
	for _, tt := range tests {
		x := &XMLSpell{
			Name:       tt.fields.Name,
			Level:      tt.fields.Level,
			School:     tt.fields.School,
			Ritual:     tt.fields.Ritual,
			Time:       tt.fields.Time,
			Range:      tt.fields.Range,
			Components: tt.fields.Components,
			Duration:   tt.fields.Duration,
			Classes:    tt.fields.Classes,
			Texts:      tt.fields.Texts,
		}
		if got := x.parseComponents(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. XMLSpell.parseComponents() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
