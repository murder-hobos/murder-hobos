package xmlspellparse

import (
	"database/sql"
	"encoding/xml"
	"strings"
)

var (
	schools = map[string]string{
		"A":  "Abjuration",
		"C":  "Conjuration",
		"D":  "Divination",
		"EN": "Enchantment",
		"EV": "Evocation",
		"I":  "Illusion",
		"N":  "Necromancy",
		"T":  "Transmutation",
	}
)

// Compendium represents a <compendium> element
type Compendium struct {
	XMLName xml.Name `xml:"compendium"`
	Spell   []Spell  `xml:"spell"`
}

// Spell is a giant catchall representing both a spell from the xml file and
// an adapted version to our database.
type Spell struct {
	Name          string   `xml:"name"db:"name"`
	Level         string   `xml:"level"db:"level"`
	SchoolAbbrv   string   `xml:"school"`
	RitualStr     string   `xml:"ritual"`
	RitualBool    bool     `db:"ritual"`
	Time          string   `xml:"time"db:"cast_time"`
	Range         string   `xml:"range"db:"range"`
	ComponentsStr string   `xml:"components"`
	Duration      string   `xml:"duration"db:"duration"`
	ClassesStr    string   `xml:"classes"`
	Texts         []string `xml:"text"`
	Concentration bool     `db:"concentration"`
	Components    *Components
	SourceText    int `db:"source_id"`
}

// Init finishes initalizing the Spell, filling out derived DB fields from
// the provided xml fields
func (s *Spell) Init() {
	s.Components = parseComponents(s.ComponentsStr)
	s.Concentration = strings.Contains(s.Duration, "Concentration")
	s.RitualBool = s.RitualStr != ""

}

// Class represents our database Class table
type Class struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	BaseClass int    `db:"base_class_id"`
}

// Components is needed because the xml file has everything on one darn line
type Components struct {
	Verb    bool           `db:"comp_verbal"`
	Som     bool           `db:"comp_somatic"`
	Mat     bool           `db:"comp_material"`
	Matdesc sql.NullString `db:"comp_material_desc"`
}

// Capitalize a single char from a string at specified index
func capitalizeAtIndex(s string, i int) string {
	out := []rune(s)
	badstr := string(out[i])
	goodstr := strings.ToUpper(badstr)
	goodrune := []rune(goodstr)
	out[i] = goodrune[0]
	return string(out)
}

// toNullString converts a regular string to a sql.NullString
func toNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// Ugly situational parser
// Parses info from a string and returns a Components struct
func parseComponents(s string) *Components {
	var verb, som, mat bool
	var matdesc sql.NullString

	// really taking advantage of the fact that spell descriptions are all lower case
	verb = strings.Contains(s, "V")
	som = strings.Contains(s, "S")
	mat = strings.Contains(s, "M")

	// ('s only ocurr in our domain as the beginning of the material description
	i := strings.Index(s, "(")
	if i > -1 {
		// Trim off parens
		desc := s[i+1 : len(s)-1]

		// Descriptions are all lower case, make them look prettier
		// by capitalizing the first letter
		cdesc := capitalizeAtIndex(desc, 0)
		matdesc = toNullString(cdesc)
	}

	return &Components{verb, som, mat, matdesc}
}
