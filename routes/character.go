package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/murder-hobos/murder-hobos/model"
	"github.com/murder-hobos/murder-hobos/util"
)

func (env *Env) characterIndex(w http.ResponseWriter, r *http.Request) {
	c := r.Context().Value("Claims")
	claims := c.(Claims)

	chars, err := env.db.GetAllCharacters(claims.UID)
	if err != nil && err != model.ErrNoResult {
		errorHandler(w, r, http.StatusInternalServerError)
	}

	data := map[string]interface{}{
		"Claims":     claims,
		"Characters": chars,
	}

	if tmpl, ok := env.tmpls["characters.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

}

// Information about specific character
func (env *Env) characterDetails(w http.ResponseWriter, r *http.Request) {
	c := r.Context().Value("Claims")
	claims := c.(Claims)
	name := mux.Vars(r)["characterName"]

	char := &model.Character{}
	c, err := env.db.GetCharacterByName(claims.UID, name)
	if err != nil {
		log.Printf("Error getting Character with name: %s\n", name)
		log.Printf(err.Error())
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Claims":    claims,
		"Character": char,
	}

	if tmpl, ok := env.tmpls["character-details.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Printf("Error loading template for class-details\n")
		return
	}
}

func (env *Env) newCharacterIndex(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value("Claims").(Claims)

	data := map[string]interface{}{
		"Claims": claims,
	}

	if tmpl, ok := env.tmpls["character-creator.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
		log.Println("EXECUTED")
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Printf("Error loading template for character-creator\n")
		return
	}
}

func (env *Env) newCharacterProcess(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("Claims").(Claims)

	name := r.PostFormValue("name")
	//class := r.PostFormValue("class")
	//level := r.PostFormValue("level")
	race := r.PostFormValue("race")
	a, _ := strconv.Atoi(r.PostFormValue("abilityMod"))
	p, _ := strconv.Atoi(r.PostFormValue("profBonus"))
	ability := util.ToNullInt64(a)
	proficiency := util.ToNullInt64(p)

	char := &model.Character{
		Name:                 name,
		Race:                 race,
		SpellAbilityModifier: ability,
		ProficienyBonus:      proficiency,
		UserID:               claims.UID,
	}

	if _, err := env.db.CreateCharacter(claims.UID, *char); err != nil {
		log.Printf("CreateCharacter: %s\n", err.Error())
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	//	if _, err := env.db.SetCharacterLevel(charID, className, level int); err != nil {
	//		errorHandler(w,r,http.StatusInternalServerError)
	//	}
	r.Method = "GET"
	http.Redirect(w, r, "/user/character", http.StatusFound)
}
