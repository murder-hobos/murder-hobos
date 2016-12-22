package initDb

import (
	"bytes"
	"database/sql"
	"errors"
	"log"
	"strings"

	"html"

	"github.com/murder-hobos/murder-hobos/model"
	"github.com/murder-hobos/murder-hobos/util"
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

const (
	// PHBid is the Player's Handbook id in our db
	PHBid = 1
	// EEid is the Elemental Evil id in our db
	EEid = 2
	// SCAGid is the Sword Coast Adventurer's guide id in our db
	SCAGid = 3
)

// XMLSpell represents a <spell> element from our xml file
//
// Struct tags in Go (the stuff between the `'s) are used commonly
// by encoding packages. Here we're telling the xml package how we
// want it to parse into our struct. Each time the xml parser encounters
// an xml element, it looks for a struct tag in our struct that matches
// that elements name. If it finds one, it assigns the value from that
// element to our struct field.
type XMLSpell struct {
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

// Compendium represents our top level <compendium> element
type Compendium struct {
	XMLSpells []XMLSpell `xml:"spell"`
}

// Components is needed because the xml file has everything on one darn line
type components struct {
	Verb    bool
	Som     bool
	Mat     bool
	Matdesc sql.NullString
}

// ToDbSpell parses the data from `x` into a new Spell object
// which it returns, along with an error. In the event of an error,
// a zero-valued Spell is returned.
func (x *XMLSpell) ToDbSpell() (model.Spell, error) {

	// vars we need to do a little work for
	// to convert
	var school, desc string
	var concentration, ritual bool
	var comps components

	sourceID := PHBid //default to phb

	if strings.Contains(x.Name, "(EE)") {
		sourceID = EEid
	}
	if strings.Contains(x.Name, "(SCAG)") {
		sourceID = SCAGid
	}

	// We want the long version, not the abbreviation
	if s, ok := schools[x.School]; ok {
		school = s
	} else {
		return model.Spell{}, errors.New("Not in schools map")
	}

	var b bytes.Buffer

	// The texts slice holds a slice of strings representing the spell
	// description from the xml file. <text/> elements are used in the file
	// to create newlines, here we replace them with <br/> to render
	// correctly in html.
	for _, text := range x.Texts {
		if text == "" {
			b.Write([]byte("\n\n"))
		}
		if text != "" {
			b.Write([]byte(html.EscapeString(text)))
		}
		// This is dirty, but the file doesn't have a field
		// for concentation, only way to find it is to see
		// if the description mentions it.
		if strings.Contains(text, "concentration") {
			concentration = true
		}
	}
	desc = b.String()

	comps = x.parseComponents()
	// In the file, ritual will be either "" or "YES"
	ritual = strings.Compare(x.Ritual, "YES") == 0

	d := model.Spell{
		Name:          trimSourceFromName(x.Name),
		Level:         x.Level,
		School:        school,
		CastTime:      x.Time,
		Duration:      x.Duration,
		Range:         x.Range,
		Verbal:        comps.Verb,
		Somatic:       comps.Som,
		Material:      comps.Mat,
		MaterialDesc:  comps.Matdesc,
		Concentration: concentration,
		Ritual:        ritual,
		Description:   desc,
		SourceID:      sourceID,
	}

	return d, nil
}

// ParseClasses converts the XMLSpell's string of comma seperated
// classes into a slice of Class objects fully initialized with
// ID and BaseClass values, ready to be inserted into our db.
func (x *XMLSpell) ParseClasses() ([]model.Class, bool) {
	cs := []model.Class{}
	split := strings.Split(x.Classes, ", ")
	for _, s := range split {
		// here Classes is a map found in classes.go
		// not in this file because it's long and ugly
		if c, ok := model.Classes[s]; ok {
			cs = append(cs, c)
		} else {
			return []model.Class{}, false
		}
	}
	return cs, true
}

func trimSourceFromName(name string) string {
	s := strings.NewReplacer(" (EE)", "", " (SCAG)", "")
	return s.Replace(name)
}

// parseComponents parses the information in the xml file's Components
// string into a Components struct literal
func (x *XMLSpell) parseComponents() components {
	var verb, som, mat bool
	var matdesc sql.NullString

	// really taking advantage of the fact that spell descriptions are all
	// lower case
	verb = strings.Contains(x.Components, "V")
	som = strings.Contains(x.Components, "S")
	mat = strings.Contains(x.Components, "M")

	// ('s only ocurr in our domain as the beginning of the material description
	// Index returns -1 if not present
	i := strings.Index(x.Components, "(")
	if i > -1 {
		// extract "inside parens" from "text text (inside parens)"
		desc := x.Components[i+1 : len(x.Components)-1]

		// Descriptions are all lower case, make them look prettier
		// by capitalizing the first letter
		cdesc, ok := util.CapitalizeAtIndex(desc, 0)
		if !ok {
			log.Printf("Error capitalizing %v at index %d\n", desc, 0)
		}
		matdesc = util.ToNullString(cdesc)
	}

	return components{
		Verb:    verb,
		Som:     som,
		Mat:     mat,
		Matdesc: matdesc,
	}
}
