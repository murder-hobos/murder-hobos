package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/murder-hobos/murder-hobos/model"
	"github.com/murder-hobos/murder-hobos/util"
)

func (env *Env) userSpellIndex(w http.ResponseWriter, r *http.Request) {
	c := r.Context().Value("Claims")
	claims := c.(Claims)

	spells, err := env.db.GetAllUserSpells(claims.UID)
	if err != nil && err != model.ErrNoResult {
		errorHandler(w, r, http.StatusInternalServerError)
	}

	data := map[string]interface{}{
		"Claims": claims,
		"Spells": spells,
	}

	if tmpl, ok := env.tmpls["user-spells.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func (env *Env) userSpellFilter(w http.ResponseWriter, r *http.Request) {
	c := r.Context().Value("Claims")
	claims := c.(Claims)

	level := r.FormValue("level")
	school := r.FormValue("school")

	spells, err := env.db.FilterUserSpells(claims.UID, level, school)
	if err != nil {
		if err == model.ErrNoResult {
			// do nothing, just show no results on page (already in template)
		} else { // something happened
			log.Printf("routes - userSpells: Error filtering cannon spells: %s\n", err.Error())
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
	}

	data := map[string]interface{}{
		"Spells": spells,
		"Claims": claims,
	}

	if tmpl, ok := env.tmpls["user-spells.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
	}
}

func (env *Env) userSpellSearch(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("Claims").(Claims)
	name := r.FormValue("name")

	spells, err := env.db.SearchUserSpells(claims.UID, name)
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

	if tmpl, ok := env.tmpls["user-spells.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
	}
}

// Show information about a single spell
func (env *Env) userSpellDetails(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("Claims").(Claims)
	name := mux.Vars(r)["spellName"]

	spell, err := env.db.GetUserSpellByName(claims.UID, name)
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
		"IsUser":  true,
	}

	if tmpl, ok := env.tmpls["spell-details.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Printf("Error loading template for spell-details\n")
		return
	}
}

func (env *Env) newSpellProcess(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("Claims").(Claims)

	r.ParseForm()
	name := r.FormValue("name")
	school := r.FormValue("school")
	level := r.FormValue("level")
	castTime := r.FormValue("castTime")
	duration := r.FormValue("duration")
	ran := r.FormValue("range")
	verbal, _ := strconv.ParseBool(r.FormValue("verbal"))
	somatic, _ := strconv.ParseBool(r.FormValue("somatic"))
	material, _ := strconv.ParseBool(r.FormValue("material"))
	materialDesc := util.ToNullString(r.FormValue("materialDesc"))
	conc, _ := strconv.ParseBool(r.FormValue("concentration"))
	ritual, _ := strconv.ParseBool(r.FormValue("ritual"))
	desc := r.FormValue("spellDesc")
	sourceID := claims.UID

	spell := &model.Spell{
		0,
		name,
		level,
		school,
		castTime,
		duration,
		ran,
		verbal,
		somatic,
		material,
		materialDesc,
		conc,
		ritual,
		desc,
		sourceID,
	}
	if _, err := env.db.CreateSpell(claims.UID, *spell); err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	http.Redirect(w, r, "/user/spell", http.StatusFound)
}

func (env *Env) newSpellIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("newSpellIndex")
	claims, _ := r.Context().Value("Claims").(Claims)

	classes, err := env.db.GetAllClasses()
	if err != nil {
		log.Println(err.Error())
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Claims":  claims,
		"Classes": classes,
	}

	if tmpl, ok := env.tmpls["spell-creator.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
		log.Println("EXECUTED")
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Printf("Error loading template for spell-creator\n")
		return
	}
}

func (env *Env) userProfileIndex(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value("Claims").(Claims)

	data := map[string]interface{}{
		"Claims": claims,
	}

	if tmpl, ok := env.tmpls["profile.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Printf("Error loading template for spell-creator\n")
		return
	}

}

func (env *Env) userSpellDelete(w http.ResponseWriter, r *http.Request) {

	spellID := mux.Vars(r)["userSpellID"]

	if _, err := env.db.DeleteSpell(spellID); err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	http.Redirect(w, r, "/user/spell", http.StatusFound)
}
