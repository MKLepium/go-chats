package httpserver

import (
	"html/template"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}

func chat(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("internal/http/chat.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
