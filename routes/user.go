package routes

import (
	"html"
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

	name := r.PostFormValue("name")
	school := r.PostFormValue("school")
	level := r.PostFormValue("level")
	castTime := r.PostFormValue("castTime")
	duration := r.PostFormValue("duration")
	ran := r.PostFormValue("range")
	verbal := r.PostFormValue("verbal") != ""
	somatic := r.PostFormValue("somatic") != ""
	material := r.PostFormValue("material") != ""
	materialDesc := util.ToNullString(r.PostFormValue("materialDesc"))
	conc := r.PostFormValue("concentration") != ""
	ritual := r.PostFormValue("ritual") != ""
	desc := html.EscapeString(r.PostFormValue("spellDesc"))
	sourceID := claims.UID
	spell := &model.Spell{
		ID:            0,
		Name:          name,
		Level:         level,
		School:        school,
		CastTime:      castTime,
		Duration:      duration,
		Range:         ran,
		Verbal:        verbal,
		Somatic:       somatic,
		Material:      material,
		MaterialDesc:  materialDesc,
		Concentration: conc,
		Ritual:        ritual,
		Description:   desc,
		SourceID:      sourceID,
	}

	if _, err := env.db.CreateSpell(claims.UID, *spell); err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	r.Method = "GET"
	http.Redirect(w, r, "/user/spell", http.StatusFound)
}

func (env *Env) newSpellIndex(w http.ResponseWriter, r *http.Request) {
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
	claims := r.Context().Value("Claims").(Claims)
	sID := r.PostFormValue("spellID")
	spellID, err := strconv.Atoi(sID)
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest)
		log.Printf("userSpellDelete: Error converting string to int")
		return
	}

	if err := env.db.DeleteSpell(claims.UID, spellID); err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	r.Method = "GET"
	http.Redirect(w, r, "/user/spell", http.StatusFound)
}
