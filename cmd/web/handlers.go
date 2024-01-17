package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/arian-nj/snippetbox/internals/models"
)

func (app *application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}
	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }
	// files := []string{
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/pages/home.tmpl",
	// }

	// tp, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }
	// err = tp.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// }
}

func (app *application) SnippetView(w http.ResponseWriter, r *http.Request) {
	id_par := r.URL.Query().Get("id")
	id_int, err := strconv.Atoi(id_par)
	if err != nil || id_int < 1 { // invalid id
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id_int)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) SnippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n–Kobayashi Issa"
	expires := 7
	inserted_id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", inserted_id), http.StatusSeeOther)
}
