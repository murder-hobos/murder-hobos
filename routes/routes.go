package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	gmux "github.com/gorilla/mux"
)

var (
	// package global holding all of our templates in our templates dir
	tmpls map[string]*template.Template
)

func init() {
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
		fmt.Println(files)
		tmpls[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
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
	r.PathPrefix("/static").Handler(staticHandler)
	return r
}

// hello world
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello murder-hobos!")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Murder Hobos", Body: "Welcome to Murder Hobos!"}
	err := renderTemplate(w, "index", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// serves static files (css/js)
var staticHandler = http.StripPrefix("/static", http.FileServer(http.Dir("./static")))
