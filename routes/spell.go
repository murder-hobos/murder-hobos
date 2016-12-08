package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/murder-hobos/murder-hobos/model"
)

func (env *Env) spellIndex(w http.ResponseWriter, r *http.Request) {
	spells, err := env.db.GetAllCannonSpells()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if tmpl, ok := env.tmpls["spells.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", spells)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func (env *Env) spellFilter(w http.ResponseWriter, r *http.Request) {
	level := r.FormValue("level")
	school := r.FormValue("school")

	spells, err := env.db.FilterCannonSpells(level, school)
	if err != nil {
		if err == model.ErrNoResult {
			// do nothing, just show no results on page (already in template)
		} else { // something happened
			log.Printf("routes - cannonSpells: Error filtering cannon spells: %s\n", err.Error())
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
	}

	if tmpl, ok := env.tmpls["spells.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", spells)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
	}
}

func (env *Env) spellSearch(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")

	spells, err := env.db.SearchCannonSpells(name)
	if err != nil {
		if err == model.ErrNoResult {
			// do nothing, just show no results on page (already in template)
		} else { // something happened
			log.Printf("routes - cannonSpells: Error filtering cannon spells: %s\n", err.Error())
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
	}

	if tmpl, ok := env.tmpls["spells.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", spells)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
	}
}

// Show information about a single spell
func (env *Env) spellDetails(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["spellName"]

	spell, err := env.db.GetCannonSpellByName(name)
	if err != nil {
		log.Printf("Error getting spell by name: %s\n", name)
		log.Printf(err.Error())
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	classes, err := env.db.GetSpellClasses(spell.ID)
	// we shouldn't have an error at this point, we should have a spell
	if err != nil {
		log.Printf("Error getting spell classes with id %d\n", spell.ID)
		log.Println(err.Error())
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if tmpl, ok := env.tmpls["spell-details.html"]; ok {
		data := struct {
			Spell   *model.Spell
			Classes *[]model.Class
		}{
			spell,
			classes,
		}
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Printf("Error loading template for spell-details\n")
		return
	}
}
