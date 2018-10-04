package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/robots.txt", robots)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		renderTemplate(w, "index", r)
		return
	}
	http.Error(w, "404 Not Found", http.StatusNotFound)
}

func robots(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("robots.txt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%s", content)
}

func renderTemplate(w http.ResponseWriter, fname string, r *http.Request) {
	t, err := template.ParseFiles("templates/base.html", "templates/"+fname+".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "base", "")
}
