package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type application struct {
	errorlog *log.Logger
	infolog  *log.Logger
}

func (app *application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	files := []string{
		"./ui/html/partials/nav.tmpl",
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	tp, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	err = tp.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) SnippetView(w http.ResponseWriter, r *http.Request) {
	id_par := r.URL.Query().Get("id")
	id_int, err := strconv.Atoi(id_par)
	if err != nil || id_int < 1 { // invalid id
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "View Snippet %d", id_int)
}

func (app *application) SnippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a Snippet"))
}
