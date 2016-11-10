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
	XMLName   xml.Name    `xml:"compendium"`
	FileSpell []FileSpell `xml:"spell"`
}

// FileSpell represents a <spell> element from our xml file
type FileSpell struct {
	Name          string   `xml:"name"`
	Level         string   `xml:"level"`
	School        string   `xml:"school"`
	Ritual        string   `xml:"ritual"`
	Time          string   `xml:"time"`
	Range         string   `xml:"range"`
	ComponentsStr string   `xml:"components"`
	Duration      string   `xml:"duration"`
	ClassesStr    string   `xml:"classes"`
	Texts         []string `xml:"text"`
	Components    *Components
}

// DbSpell represents our database version of a spell
type DbSpell struct {
	Name     string `db:"name"`
	Level    string `db:"level"`
	School   string `db:"school"`
	CastTime string `db:"cast_time"`
	Duration string `db:"duration"`
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
