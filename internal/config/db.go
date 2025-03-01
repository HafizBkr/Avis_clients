package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// DB est l'instance de la base de données disponible pour toute l'application
var DB *sql.DB

// LoadEnv charge les variables d'environnement depuis le fichier .env
func LoadEnv() {
	// Trouver le chemin absolu du répertoire racine du projet
	// On suppose que le dossier `internal` est à la racine du projet
	projectRoot, err := filepath.Abs("../../")
	if err != nil {
		log.Fatalf("Erreur lors de la récupération du répertoire racine : %v", err)
	}

	// Charger le fichier .env depuis la racine du projet
	envPath := filepath.Join(projectRoot, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Println("Fichier .env non trouvé, utilisation des variables d'environnement système")
	}
}

// InitDB initialise la connexion à la base de données PostgreSQL
func InitDB() (*sql.DB, error) {
	// Charger les variables d'environnement
	LoadEnv()

	// Récupérer la chaîne de connexion complète
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Fallback sur les paramètres individuels si DATABASE_URL n'est pas défini
		dbHost := getEnv("DB_HOST", "localhost")
		dbPort := getEnv("DB_PORT", "5432")
		dbUser := getEnv("DB_USER", "postgres")
		dbPassword := getEnv("DB_PASSWORD", "postgres")
		dbName := getEnv("DB_NAME", "avis_clients_db")
		dbSSLMode := getEnv("DB_SSLMODE", "disable") // Désactiver SSL si nécessaire

		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
		)
	}

	// Connexion à PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("échec d'ouverture de la connexion : %v", err)
	}

	// Vérifier la connexion
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("échec de connexion à la base de données : %v", err)
	}

	// Assigner la connexion à la variable globale
	DB = db

	log.Println("Connexion à la base de données PostgreSQL établie")

	// Créer la table avis_clients si elle n'existe pas
	sqlTable := `
	CREATE TABLE IF NOT EXISTS avis_clients (
		id SERIAL PRIMARY KEY,
		nom VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL,
		avis TEXT NOT NULL,
		note REAL NOT NULL CHECK (note >= 0 AND note <= 5)
	);
	`
	if _, err = DB.Exec(sqlTable); err != nil {
		return nil, fmt.Errorf("erreur lors de la création de la table : %v", err)
	}

	return DB, nil
}

// getEnv récupère la valeur d'une variable d'environnement ou retourne une valeur par défaut
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// CloseDB ferme proprement la connexion à la base de données
func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Erreur lors de la fermeture de la connexion : %v", err)
		} else {
			log.Println("Connexion à la base de données fermée")
		}
	}
}
