package main

import (
	"regexp"
	"html/template"
	"os"
	"net/http"
	"log"
)

// define page as a struct with two fields representing the tile and body
type Page struct {
	Title string
	Body []byte
}

// global variable
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")


func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil{
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}


func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
	err := templates.ExecuteTemplate(w, tmpl+".html",p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


}

func viewHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil{
		http.Redirect( w, r, "/edit/"+title,http.StatusFound)
		return

	}	
	renderTemplate(w,"view",p)
}

func editHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w,"edit",p)
}

func saveHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
