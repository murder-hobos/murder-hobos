package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/murder-hobos/murder-hobos/model"
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
