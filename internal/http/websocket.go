package httpserver

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // You may want to add origin checks for security
	},
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("chatSession")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Now you can access claims.Username for the username
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("[%s] %s", claims.Username, msg)

		// Broadcast the message to all connected clients
		// In a real application, you should handle this better
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
