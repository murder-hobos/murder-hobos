package xmlspellparse

import (
	"database/sql"
	"encoding/xml"
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
	XMLName  xml.Name   `xml:"compendium"`
	XMLSpell []XMLSpell `xml:"spell"`
}

// XMLSpell represents a <spell> element from our xml file
//
// Struct tags in Go (the stuff between the `'s) are used commonly
// by encoding packages. Here we're telling the xml package how we
// want it to parse into our struct.
// For each element, the form is
// 			`xml:"element_name_in_file,which_part_of_element"`
// Where the name in the file defaults to the lowercase version of
// the struct field, and the part of element refers to whether
// the data is an attribute (attr), part of the regular character data
// between opening and closing tags (chardata), or other options we
// won't need in this situation.
//
// When we name our struct fields the same as the file elements,
// we don't need to include the first part of that struct tag.
// Since our elements don't have any attributes, only chardata,
// the xml package knows to just throw the chardata into our fields.
type XMLSpell struct {
	Name       string   `xml:"name,chardata"`
	Level      string   `xml:"level,chardata"`
	School     string   `xml:"school,chardata"`
	Ritual     string   `xml:"ritual,chardata"`
	Time       string   `xml:"time,chardata"`
	Range      string   `xml:"range,chardata"`
	Components string   `xml:"components,chardata"`
	Duration   string   `xml:"duration,chardata"`
	Classes    string   `xml:"classes,chardata"`
	Texts      []string `xml:"text,chardata"`
}

// UnmarshalXML unmarshalls xml data from a Decoder xml parser into
// a Compendium struct
func (c *Compendium) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data *Compendium
	err := d.DecodeElement(data, &start)
	if err != nil {
		c = data
	}
	return err
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
		cdesc := capitalizeAtIndex(desc, 0)
		matdesc = toNullString(cdesc)
	}

	return &Components{verb, som, mat, matdesc}
}
