// xmlparse parses spells from our specific xml file into our specific
// database. It is assumed that before running this program, the database
// has been initialized with all data for Classes and Users.
//
// Database connection is handled in this package's config.json file
// Usage:
// 		xmlparse -f xmlfile
package main

import (
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/jaden-young/murder-hobos/util/xmlspellparse"
)

var xmlFile string
var configFile string

func init() {
	flag.StringVar(&xmlFile, "f", "", "xml file to parse")
}

func main() {
	flag.Parse()

	db, err := xmlspellparse.InitDB()
	defer db.Close()
	if err != nil {
		log.Fatalln(err)
	}

	// Have to be silly about this because range is a reserved word
	insertSpell, err := db.Prepare(`
		INSERT INTO Spell (name, level, school, cast_time, duration, 
	` + "`range`" + `, comp_verbal, comp_somatic, comp_material, material_desc, concentration, ritual, description, source_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		log.Fatalln(err)
	}

	insertClassSpells, err := db.Prepare(`
		INSERT INTO ClassSpells (spell_id, class_id) VALUES (?, ?);`)
	if err != nil {
		log.Fatalln(err)
	}

	// Find our file
	xFile, err := os.Open(xmlFile)
	defer xFile.Close()
	if err != nil {
		log.Fatalln(err)
	}

	// Read file into memory
	xBytes, err := ioutil.ReadAll(xFile)
	if err != nil {
		log.Fatalln(err)
	}
	xFile.Close()

	var c xmlspellparse.Compendium
	xml.Unmarshal(xBytes, &c)

	// for each spell in our xml file
	for _, xmlSpell := range c.XMLSpells {
		s, err := xmlSpell.ToDbSpell()
		if err != nil {
			log.Fatal("Error converting to db spell")
		}

		// Insert into spell table
		result, err := insertSpell.Exec(
			s.Name, s.Level, s.School, s.CastTime, s.Duration, s.Range, s.Components.Verb,
			s.Components.Som, s.Components.Mat, s.Components.Matdesc, s.Concentration,
			s.Ritual, s.Description, s.SourceID)
		if err != nil {
			log.Fatalln(err)
		}

		// Remember which spell we inserted so we can insert into ClassSpells
		spellID, err := result.LastInsertId()
		if err != nil {
			log.Fatalln(err)
		}

		// Insert into ClassSpells table
		if classes, ok := xmlSpell.ParseClasses(); ok {
			for _, class := range classes {
				if _, err := insertClassSpells.Exec(spellID, class.ID); err != nil {
					log.Fatalln(err)
				}
			}
		} else {
			log.Fatalf("Error parsing classes from %v", xmlSpell)
		}
	}
}
