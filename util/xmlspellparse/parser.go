package xmlspellparse

import (
	"database/sql"
	"encoding/xml"
	"log"
	"strings"
)

var (
	// Schools are abreviated in xml file, we want full text
	// in our db
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

// Compendium represents our top level <compendium> element
type Compendium struct {
	XMLSpell []XMLSpell `xml:"spell"`
}

// XMLSpell represents a <spell> element from our xml file
//
// Struct tags in Go (the stuff between the `'s) are used commonly
// by encoding packages. Here we're telling the xml package how we
// want it to parse into our struct. Each time the xml parser encounters
// an xml element, it looks for a struct tag in our struct that matches
// that elements name. If it finds one, it assigns the value from that
// element to our struct field.
type XMLSpell struct {
	XMLName    xml.Name `xml:"spell"`
	Name       string   `xml:"name"`
	Level      string   `xml:"level"`
	School     string   `xml:"school"`
	Ritual     string   `xml:"ritual"`
	Time       string   `xml:"time"`
	Range      string   `xml:"range"`
	Components string   `xml:"components"`
	Duration   string   `xml:"duration"`
	Classes    string   `xml:"classes"`
	Texts      []string `xml:"text"`
}

// DbSpell represents our database version of a spell
type DbSpell struct {
	Name          string `db:"name"`
	Level         string `db:"level"`
	School        string `db:"school"`
	CastTime      string `db:"cast_time"`
	Duration      string `db:"duration"`
	Range         string `db:"range"`
	Components    *Components
	MaterialDesc  string `db:"material_desc"`
	Concentration bool   `db:"concentration"`
	Ritual        bool   `db:"ritual"`
	Description   string `db:"description"`
	SourceID      int    `db:"source_id"`
}

// Components is needed because the xml file has everything on one darn line
type Components struct {
	Verb    bool           `db:"comp_verbal"`
	Som     bool           `db:"comp_somatic"`
	Mat     bool           `db:"comp_material"`
	Matdesc sql.NullString `db:"material_desc"`
}

// Class represents our database Class table
type Class struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	BaseClass int    `db:"base_class_id"`
}

// Capitalize a single char from a string at specified index
// If an error is encountered, normally index being out of range,
// ok will be set to false and the original string returned unaltered
func capitalizeAtIndex(s string, i int) (string, bool) {
	if i < 0 || i > len(s)-1 {
		return s, false
	}
	out := []rune(s)
	badstr := string(out[i])
	goodstr := strings.ToUpper(badstr)
	goodrune := []rune(goodstr)
	out[i] = goodrune[0]
	return string(out), true
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

	// really taking advantage of the fact that spell descriptions are all
	// lower case
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
		cdesc, ok := capitalizeAtIndex(desc, 0)
		if !ok {
			log.Printf("Error capitalizing %v at index %d\n", desc, 0)
		}
		matdesc = toNullString(cdesc)
	}

	return &Components{verb, som, mat, matdesc}
}
