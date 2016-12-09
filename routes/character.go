package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/murder-hobos/murder-hobos/model"
)

// List characters
func (env *Env) characterIndex(w http.ResponseWriter, r *http.Request) {
	var userID int

	//we will need to make this specific to userID soon on an account by account basis
	characters, err := env.db.GetAllCharacters(userID)
	if err != nil {
		if err.Error() == "empty slice passed to 'in' query" || err == model.ErrNoResult {
			// do nothing, just show no results on page (already in template)
		} else { // something happened
			log.Println(err.Error())
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
	}

	if tmpl, ok := env.tmpls["characters.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", characters)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
	}
}

// Information about specific character
func (env *Env) characterDetails(w http.ResponseWriter, r *http.Request) {
	var userID int
	name := mux.Vars(r)["characterName"]

	c := &model.Character{}
	c, err := env.db.GetCharacterByName(userID, name)
	if err != nil {
		log.Printf("Error getting Character with name: %s\n", name)
		log.Printf(err.Error())
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	data := struct {
		Character *model.Character
	}{
		c,
	}
	if tmpl, ok := env.tmpls["character-details.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Printf("Error loading template for class-details\n")
		return
	}
}
