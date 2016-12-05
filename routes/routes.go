package routes

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/murder-hobos/murder-hobos/model"
)

var (
	// Package global that holds a map of our templates.
	tmpls map[string]*template.Template
)

const (
	sessionKey = "murder-hobos"
)

// Env is a struct that defines an enviornment for server request handling.
// It allows us to specify different combinations of datastores, templates,
// and session stores
type Env struct {
	db    model.Datastore
	tmpls map[string]*template.Template
	store *sessions.CookieStore
}

func init() {
	// setup template map
	if tmpls == nil {
		tmpls = make(map[string]*template.Template)
	}

	templatesDir := "templates/"

	layouts, err := filepath.Glob(templatesDir + "layouts/*.html")
	if err != nil {
		log.Fatal(err)
	}

	includes, err := filepath.Glob(templatesDir + "includes/*.html")
	if err != nil {
		log.Fatal(err)
	}

	for _, layout := range layouts {
		files := append(includes, layout)
		tmpls[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
	}
}

// New returns an initialized gorilla *mux.Router with all of our routes
// Panics if unable to connect to datastore with given dsn
// (don't want the server to start without database access)
func New(dsn string) *mux.Router {
	store := sessions.NewCookieStore([]byte("super-secret-key-that-is-totally-secure"))

	db, err := model.NewDB(dsn)
	if err != nil {
		panic(err)
	}
	env := &Env{db, tmpls, store}

	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/class/{className}", env.classDetailsHandler)
	r.HandleFunc("/classes", env.classesHandler)
	r.HandleFunc("/spell/{spellName}", env.spellDetailsHandler)
	r.HandleFunc("/spells", env.spellsHandler)
	r.PathPrefix("/static").HandlerFunc(staticHandler)
	return r
}

// Index doesn't really do much for now
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if tmpl, ok := tmpls["index.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", nil)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
	}
}

// List all spells. We really should chache this eventually
// instead of hitting the db everytime
func (env *Env) spellsHandler(w http.ResponseWriter, r *http.Request) {
	var userID int
	includeCannon := true // want to default to true, not false

	if i, ok := env.getIntFromSession(r, "userID"); ok {
		userID = i
	}
	if b, ok := env.getBoolFromSession(r, "includeCannon"); ok {
		includeCannon = b
	}

	spells, err := env.db.GetAllSpells(userID, includeCannon)
	if err != nil {
		if err.Error() == "empty slice passed to 'in' query" || err == model.ErrNoResult {
			// do nothing, just show no results on page (already in template)
		} else { // something happened
			log.Println(err.Error())
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
func (env *Env) spellDetailsHandler(w http.ResponseWriter, r *http.Request) {
	var userID int
	includeCannon := true // want to default to true, not false

	if i, ok := env.getIntFromSession(r, "userID"); ok {
		userID = i
	}
	if b, ok := env.getBoolFromSession(r, "includeCannon"); ok {
		includeCannon = b
	}

	name := mux.Vars(r)["spellName"]

	s, err := env.db.GetSpellByName(name, userID, includeCannon)
	if err != nil {
		log.Printf("Error getting spell by name: %s, userID: %d, isCannon: %t\n", name, s.ID, true)
		log.Printf(err.Error())
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	cs, err := env.db.GetSpellClasses(s.ID)
	// we shouldn't have an error at this point, we should have a spell
	if err != nil {
		log.Printf("Error getting spell classes with id %d\n", s.ID)
		log.Println(err.Error())
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if tmpl, ok := env.tmpls["spell-details.html"]; ok {
		data := struct {
			Spell   *model.Spell
			Classes *[]model.Class
		}{
			s,
			cs,
		}
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Printf("Error loading template for spell-details\n")
		return
	}
}

// lists all classes
func (env *Env) classesHandler(w http.ResponseWriter, r *http.Request) {
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

//list a single class and all spells available to that class
func (env *Env) classDetailsHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["className"]

	s, err := env.db.GetClassByName(name)
	if err != nil {
		log.Printf("Error getting Class by name: %s, classID: %d\n", name, s.ID)
		log.Printf(err.Error())
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	cs, err := env.db.GetClassSpells(s.ID)
	if err != nil {
		log.Println("Class-detail handler" + err.Error())
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if tmpl, ok := env.tmpls["class-details.html"]; ok {
		tmpls.ExecuteTemplate(w, "base", cs)
	} else {
		errorhandler(w, r, http.StatusInternalServerError)
		log.Printf("Error loading template for class-details\n")
		return
	}
}

// serve static (js/css) files
func staticHandler(w http.ResponseWriter, r *http.Request) {
	// Don't want to list directories
	if strings.HasSuffix(r.URL.Path, "/") {
		http.Error(w, "File not found", http.StatusBadRequest)
	}
	http.ServeFile(w, r, r.URL.Path[1:])
}

// Custom stuff for errors
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	var title, message string
	var tmpl *template.Template

	if t, ok := tmpls["error.html"]; ok {
		tmpl = t
	} else {
		http.Error(w, "Server's busted.", http.StatusInternalServerError)
	}

	if status == http.StatusNotFound {
		title = "Not Found"
		message = "Whoops! We can't find that!"
	}

	if status == http.StatusInternalServerError {
		title = "Server Error"
		message = "Our server is having issues. >:("
	}

	vars := map[string]string{"Title": title, "Message": message}
	tmpl.ExecuteTemplate(w, "base", vars)
}

// litle utils
func (env *Env) getStringFromSession(r *http.Request, key string) (string, bool) {
	sess, err := env.store.Get(r, sessionKey)
	if err != nil {
		return "", false
	}
	val, ok := sess.Values[key]
	if !ok {
		return "", false
	}
	if s, ok := val.(string); ok {
		return s, true
	}
	return "", false
}

func (env *Env) getIntFromSession(r *http.Request, key string) (int, bool) {
	sess, err := env.store.Get(r, sessionKey)
	if err != nil {
		return 0, false
	}
	val, ok := sess.Values[key]
	if !ok {
		return 0, false
	}
	if i, ok := val.(int); ok {
		return i, true
	}
	return 0, false
}

func (env *Env) getBoolFromSession(r *http.Request, key string) (bool, bool) {
	sess, err := env.store.Get(r, sessionKey)
	if err != nil {
		return false, false
	}
	val, ok := sess.Values[key]
	if !ok {
		return false, false
	}
	if b, ok := val.(bool); ok {
		return b, true
	}
	return false, false
}
