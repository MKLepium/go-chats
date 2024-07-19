package httpserver

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	var err error
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PW")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	for i := 0; i < 5; i++ {
		if err = db.Ping(); err != nil {
			log.Printf("Attempt %d: failed to ping: %v\n", i+1, err)
			time.Sleep(time.Second * 5)
		} else {
			break
		}
	}
}

func StartServer() {
	initDB()
	// create a new http server
	port := os.Getenv("HTTP_SERVERPORT")
	if port == "" {
		log.Default().Println("HTTP_SERVERPORT not set, using default port 8080")
		port = ":8080"
	}

	srv := &http.Server{
		Addr:         port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// register the routes
	http.HandleFunc("/", root)
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			login(w, r) // Assumes 'login' serves the HTML form
		} else if r.Method == "POST" {
			handleLogin(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/chat", chat)
	http.HandleFunc("/ws", wsEndpoint)

	// start the server
	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
