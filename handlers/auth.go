package handlers

import (
	"database/sql"
	"html/template"
	"little-task-go/db"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var userID int
		err := db.DB.QueryRow("SELECT id FROM users WHERE username = ? AND password = ?", username, password).Scan(&userID)
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		http.Redirect(w, r, "/books", http.StatusFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}
