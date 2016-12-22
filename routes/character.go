package routes

import (
	"log"
	"net/http"
	"strings"

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
	claims, ok := r.Context().Value("Claims").(Claims)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	name := mux.Vars(r)["charName"]

	char, err := env.db.GetCharacterByName(claims.UID, name)
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

	dummy := model.ClassLevelView{
		Class: model.Class{
			Name: "",
		},
		Level: 0,
	}

	data := map[string]interface{}{
		"Claims": claims,
		"Character": &model.Character{
			Levels: []model.ClassLevelView{dummy},
		},
	}

	if tmpl, ok := env.tmpls["character-creator.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Printf("Error loading template for character-creator\n")
		return
	}
}

func (env *Env) newCharacterProcess(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value("Claims").(Claims)

	name := r.PostFormValue("name")
	classes := strings.Split(r.PostFormValue("classes[]"), ", ")
	levels := strings.Split(r.PostFormValue("levels[]"), ", ")
	race := r.PostFormValue("race")
	ability := util.ToNullInt64(r.PostFormValue("abilityMod"))
	proficiency := util.ToNullInt64(r.PostFormValue("profBonus"))

	var classLevels []model.ClassLevelView
	for i, c := range classes {
		class := model.Classes[c]
		level, err := strconv.Atoi(levels[i])
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		classLevels = append(classLevels, model.ClassLevelView{
			Class: class,
			Level: level,
		})
	}

	char := &model.Character{
		Name:                 name,
		Race:                 race,
		SpellAbilityModifier: ability,
		ProficiencyBonus:     proficiency,
		UserID:               claims.UID,
		Levels:               classLevels,
	}

	if _, err := env.db.InsertCharacter(char); err != nil {
		log.Printf("CreateCharacter: %s\n", err.Error())
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	r.Method = "GET"
	http.Redirect(w, r, "/user/character", http.StatusFound)
}

func (env *Env) editCharacterIndex(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value("Claims").(Claims)
	name := mux.Vars(r)["charName"]

	char, err := env.db.GetCharacterByName(claims.UID, name)
	if err != nil {
		errorHandler(w, r, http.StatusNotFound)
		log.Println("Error getting char by name")
		log.Println(err.Error())
		return
	}

	data := map[string]interface{}{
		"Claims":    claims,
		"Character": char,
	}

	if tmpl, ok := env.tmpls["character-creator.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Printf("Error loading template for character-creator\n")
		return
	}

}
