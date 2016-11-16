package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	gmux "github.com/gorilla/mux"
)

var (
	// Package global that holds a map of our templates.
	//
	tmpls map[string]*template.Template
)

// package initialization, sets up our tmpls map
func init() {
	log.Println("Initializing template map")
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
		log.Printf("Added to tmpl map: %s\n", filepath.Base(layout))
	}
}

func renderTemplate(w http.ResponseWriter, name string, p *Page) error {
	tmpl, ok := tmpls[name+".html"]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.ExecuteTemplate(w, "base", p)
}

// Page represents a basic page
type Page struct {
	Title string
	Body  string
}

// New returns an initialized gorilla *mux.Router with all of our routes
func New() *gmux.Router {
	r := gmux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.PathPrefix("/static").HandlerFunc(staticHandler)
	return r
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Murder Hobos", Body: "Welcome to Murder Hobos!"}
	err := renderTemplate(w, "index", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	// Disable listing directories
	if strings.HasSuffix(r.URL.Path, "/") {
		http.Error(w, "File not found", http.StatusBadRequest)
	}
	http.ServeFile(w, r, r.URL.Path[1:])
}
