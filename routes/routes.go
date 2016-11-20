package routes

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"os"

	"github.com/go-sql-driver/mysql"
	mux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jaden-young/murder-hobos/model"
	"github.com/jmoiron/sqlx"
)

var (
	// Package global that holds a map of our templates.
	tmpls map[string]*template.Template

	// DB connection
	db *model.DB

	// Session store
	store = sessions.NewCookieStore([]byte("super-secret-key-that-is-totally-secure"))
)

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

	dbconfig := mysql.Config{
		User:            os.Getenv("MYSQL_USER"),
		Passwd:          os.Getenv("MYSQL_PASS"),
		DBName:          os.Getenv("MYSQL_DB_NAME"),
		Net:             "tcp",
		Addr:            os.Getenv("MYSQL_ADDR"),
		MultiStatements: false,
	}

	conn := sqlx.MustConnect("mysql", dbconfig.FormatDSN())
	db = &model.DB{DB: conn}

}

// New returns an initialized gorilla *mux.Router with all of our routes
func New() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/spell/{spellName}", db.WithDB(withSourceIDs(spellDetailsHandler)))
	r.HandleFunc("/spells", db.WithDB(withSourceIDs(spellsHandler)))
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

// /spells
// List all spells. We really should chache this eventually
// instead of hitting the db everytime
func spellsHandler(w http.ResponseWriter, r *http.Request) {
	sourceIDs := r.Context().Value("sourceIDs").([]string)
	db := r.Context().Value("db").(*model.DB)

	spells, ok := db.GetAllSpells(sourceIDs)
	if !ok {
		log.Printf("Error getting all spells\n")
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if tmpl, ok := tmpls["spells.html"]; ok {
		tmpl.ExecuteTemplate(w, "base", spells)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
	}
}

// Show information about a single spell
func spellDetailsHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["spellName"]
	sourceIDs := r.Context().Value("sourceIDs").([]string)
	db := r.Context().Value("db").(*model.DB)
	var spell *model.Spell
	var classes *[]model.Class

	if s, ok := db.GetSpellByName(name, sourceIDs); ok {
		spell = s
	} else {
		errorHandler(w, r, http.StatusNotFound)
		log.Printf("Couldn't find spell with name: %s and ids: %s\n", name, sourceIDs)
		return
	}

	if cs, err := db.GetSpellClasses(spell.ID); err == nil {
		classes = cs
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Printf("Error getting spell classes: %s\n", err.Error())
		return
	}

	if tmpl, ok := tmpls["spell-details.html"]; ok {
		data := struct {
			Spell   *model.Spell
			Classes *[]model.Class
		}{
			spell,
			classes,
		}
		tmpl.ExecuteTemplate(w, "base", data)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
		log.Printf("Error loading template for spell-details\n")
		return
	}
}
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
