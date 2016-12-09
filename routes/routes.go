package routes

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/murder-hobos/murder-hobos/model"
)

var (
	// Package global that holds a map of our templates.
	tmpls map[string]*template.Template
)

// Env is a struct that defines an enviornment for server request handling.
// It allows us to specify different combinations of datastores, templates,
type Env struct {
	db    model.Datastore
	tmpls map[string]*template.Template
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
	db, err := model.NewDB(dsn)
	if err != nil {
		panic(err)
	}
	env := &Env{db, tmpls}

	stdChain := alice.New(env.withClaims)

	r := mux.NewRouter()
	r.Handle("/spell", newSpellRouter(env))
	r.Handle("/class", newClassRouter(env))
	r.HandleFunc("/character/{characterName}", env.characterDetails)
	r.HandleFunc("/character", env.characterIndex)
	r.Handle("/", stdChain.ThenFunc(rootIndex))

	r.HandleFunc("/login", env.loginIndex).Methods("GET")
	r.HandleFunc("/login", env.loginProcess).Methods("POST")
	//r.HandleFunc("/login/register", env.loginRegister).Methods("POST")
	r.HandleFunc("/logout", env.logoutProcess)

	r.PathPrefix("/static").HandlerFunc(staticHandler)
	return r
}

// Index doesn't really do much for now
func rootIndex(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("Claims")

	data := map[string]interface{}{
		"Claims": claims,
	}

	if tmpl, ok := tmpls["index.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", data)
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
