package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/murder-hobos/murder-hobos/model"
)

// lists all classes
func (env *Env) classIndex(w http.ResponseWriter, r *http.Request) {
	cs, err := env.db.GetAllClasses()
	if err != nil {
		log.Println("Classes handler: " + err.Error())
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if tmpl, ok := env.tmpls["classes.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", cs)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Printf("Error loading template for classes\n")
		return
	}
}

// Shows a list of all spells available to a class
func (env *Env) classDetails(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["className"]

	c, err := env.db.GetClassByName(name)
	if err != nil {
		log.Printf("Error getting Class by name: %s\n", name)
		log.Printf(err.Error())
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	spells, err := env.db.GetClassSpells(c.ID)
	if err != nil {
		log.Println("Class-detail handler" + err.Error())
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	data := struct {
		Class  *model.Class
		Spells *[]model.Spell
	}{
		c,
		spells,
	}
	if tmpl, ok := env.tmpls["class-details.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Printf("Error loading template for class-details\n")
		return
	}
}
