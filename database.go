package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB() *sql.DB {

	db, err := sql.Open("sqlite", "chat.db")
	if err != nil {
		log.Fatalf("Erreur d'ouverture de la base de données : %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Erreur de connexion à SQLite : %v", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		auteur TEXT NOT NULL,
		contenu TEXT NOT NULL,
		date TEXT NOT NULL
	);
	`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Erreur lors de la création de la table messages : %v", err)
	}

	fmt.Println("📦 [BASE DE DONNÉES] SQLite initialisé avec succès et table 'messages' prête.")

	return db
}
