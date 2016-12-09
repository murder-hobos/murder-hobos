package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/murder-hobos/murder-hobos/model"
)

func newSpellRouter(env *Env) *mux.Router {
	stdChain := alice.New(env.withClaims)
	r := mux.NewRouter()

	r.HandleFunc(`/spell/{spellName:[a-zA-Z '\-\/]+}`, env.spellDetails)
	r.Handle("/spell", stdChain.ThenFunc(env.spellSearch)).Queries("name", "")
	r.Handle("/spell", stdChain.ThenFunc(env.spellFilter)).Queries("school", "")
	r.Handle("/spell", stdChain.ThenFunc(env.spellFilter)).Queries("level", "{level:[0-9]}")
	r.Handle("/spell", stdChain.ThenFunc(env.spellFilter)).Queries("school", "", "level", "{level:[0-9]}")
	r.Handle("/spell", stdChain.ThenFunc(env.spellIndex))

	return r
}

func (env *Env) spellIndex(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("Claims")

	spells, err := env.db.GetAllCannonSpells()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := map[string]interface{}{
		"Spells": spells,
		"Claims": claims,
	}

	if tmpl, ok := env.tmpls["spells.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func (env *Env) spellFilter(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("Claims")

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

	data := map[string]interface{}{
		"Spells": spells,
		"Claims": claims,
	}

	if tmpl, ok := env.tmpls["spells.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
	}
}

func (env *Env) spellSearch(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims")
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

	data := map[string]interface{}{
		"Spells": spells,
		"Claims": claims,
	}

	if tmpl, ok := env.tmpls["spells.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
	}
}

// Show information about a single spell
func (env *Env) spellDetails(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims")
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

	data := map[string]interface{}{
		"Spell":   spell,
		"Classes": classes,
		"Claims":  claims,
	}

	if tmpl, ok := env.tmpls["spell-details.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Printf("Error loading template for spell-details\n")
		return
	}
}
