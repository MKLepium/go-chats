package mysql

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
)

func AuthenticateUser(db *sql.DB, username string, password string) bool {
	query := `SELECT password_hash FROM users WHERE username = ?;`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var hashedPassword string
	err = stmt.QueryRow(username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// Username not found
			log.Printf("Username %s not found\n", username)
			return false
		}
		log.Fatal(err)
	}

	// Hash the provided password using SHA-256 to match the database
	hasher := sha256.New()
	hasher.Write([]byte(password))
	passwordSHA256 := hex.EncodeToString(hasher.Sum(nil))

	// debug prints
	log.Printf("hashedPassword: %s\n", hashedPassword)
	log.Printf("passwordSHA256: %s\n", passwordSHA256)

	// Compare the SHA-256 hashed password with the database hash
	return hashedPassword == passwordSHA256
}
