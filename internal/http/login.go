package httpserver

import (
	"html/template"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mklepium/chats/internal/mysql"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func login(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("internal/http/login.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check for an error message in the cookie
	cookie, err := r.Cookie("loginError")
	var errorMessage string
	if err == nil && cookie != nil {
		errorMessage = cookie.Value
		// Clear the cookie after reading it
		http.SetCookie(w, &http.Cookie{
			Name:     "loginError",
			Value:    "",
			Expires:  time.Unix(0, 0),
			Path:     "/",
			HttpOnly: true,
		})
	}

	data := struct {
		ErrorMessage string
	}{
		ErrorMessage: errorMessage,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if !mysql.AuthenticateUser(db, username, password) {
		http.SetCookie(w, &http.Cookie{
			Name:     "loginError",
			Value:    "Invalid username or password",
			Path:     "/",
			HttpOnly: true, // Important for security, makes the cookie inaccessible to JavaScript
		})
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Create JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set JWT token in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "chatSession",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true, // Important for security, makes the cookie inaccessible to JavaScript
		Expires:  expirationTime,
	})

	http.Redirect(w, r, "/chat", http.StatusFound)
}
