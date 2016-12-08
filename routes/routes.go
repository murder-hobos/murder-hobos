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

	r.HandleFunc("/", rootIndex)
	r.HandleFunc("/character/{characterName}", env.characterDetails)
	r.HandleFunc("/character", env.characterIndex)
	r.HandleFunc("/class/{className}", env.classDetails)
	r.HandleFunc("/class", env.classIndex)
	r.HandleFunc(`/spell/{spellName:[a-zA-Z '\-\/]+}`, env.spellDetails)

	// TODO: explicitly list schools
	// There is a better way than listing permutations. There
	// has to be.
	r.HandleFunc("/spell", env.spellSearch).Queries("name", "")
	r.HandleFunc("/spell", env.spellFilter).Queries("school", "")
	r.HandleFunc("/spell", env.spellFilter).Queries("level", "{level:[0-9]}")
	r.HandleFunc("/spell", env.spellFilter).Queries("school", "", "level", "{level:[0-9]}")
	r.HandleFunc("/spell", env.spellIndex)

	r.PathPrefix("/static").HandlerFunc(staticHandler)
	return r
}

// Index doesn't really do much for now
func rootIndex(w http.ResponseWriter, r *http.Request) {
	if tmpl, ok := tmpls["index.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", nil)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
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

func (env *Env) loginHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := env.store.Get(r, sessionKey)
	if err != nil {
		http.Error(w, "Real broke.", http.StatusInternalServerError)
		return
	}

	u, ok := sess.Values["user"]
	if !ok {

	}

	if t, ok := tmpls["login.html"]; ok {
		t.ExecuteTemplate(w, "base", nil)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
	}
}

func (env *Env) getFromSession(r *http.Request, key string, destType interface{}) (interface{}, bool) {
	sess, err := env.store.Get(r, sessionKey)
	if err != nil {
		return nil, false
	}
	val, ok := sess.Values[key]
	if !ok {
		return nil, false
	}

	switch destType := destType.(type) {
	case int:

		if v, ok := val.(t); ok {
			return v, true
		}

	}
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
