package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"api-avis-clients/internal/model"
)

// AvisRepository définit les méthodes pour accéder aux avis clients dans la base de données
type AvisRepository struct {
	DB *sql.DB
}

// NewAvisRepository crée une nouvelle instance de AvisRepository
func NewAvisRepository(db *sql.DB) *AvisRepository {
	return &AvisRepository{DB: db}
}
// CreateMultiple ajoute plusieurs avis clients à la base de données avec une seule requête SQL
func (r *AvisRepository) CreateMultiple(avisList []model.AvisClient) ([]model.AvisClient, error) {
	// Créer la requête SQL dynamique
	var placeholders []string
	var values []interface{}

	for i, avis := range avisList {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		values = append(values, avis.Nom, avis.Email, avis.Avis, avis.Note)
	}

	// Requête SQL
	sqlStatement := fmt.Sprintf(`
		INSERT INTO avis_clients (nom, email, avis, note)
		VALUES %s
		RETURNING id`, strings.Join(placeholders, ", "))

	// Exécuter la requête
	rows, err := r.DB.Query(sqlStatement, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var createdAvis []model.AvisClient
	for rows.Next() {
		var avis model.AvisClient
		if err := rows.Scan(&avis.ID); err != nil {
			return nil, err
		}
		createdAvis = append(createdAvis, avis)
	}

	return createdAvis, nil
}

// GetAll récupère tous les avis clients
func (r *AvisRepository) GetAll() ([]model.AvisClient, error) {
	rows, err := r.DB.Query("SELECT id, nom, email, avis, note FROM avis_clients ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var avisList []model.AvisClient
	for rows.Next() {
		var avis model.AvisClient
		if err := rows.Scan(&avis.ID, &avis.Nom, &avis.Email, &avis.Avis, &avis.Note); err != nil {
			return nil, err
		}
		avisList = append(avisList, avis)
	}

	return avisList, nil
}

// GetByID récupère un avis client par son ID
func (r *AvisRepository) GetByID(id int) (model.AvisClient, error) {
	var avis model.AvisClient
	err := r.DB.QueryRow("SELECT id, nom, email, avis, note FROM avis_clients WHERE id = $1", id).
		Scan(&avis.ID, &avis.Nom, &avis.Email, &avis.Avis, &avis.Note)
	if err != nil {
		if err == sql.ErrNoRows {
			return avis, errors.New("avis non trouvé")
		}
		return avis, err
	}
	return avis, nil
}

// Create ajoute un nouvel avis client
func (r *AvisRepository) Create(avis model.AvisClient) (model.AvisClient, error) {
	sqlStatement := `
	INSERT INTO avis_clients (nom, email, avis, note)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	err := r.DB.QueryRow(sqlStatement, avis.Nom, avis.Email, avis.Avis, avis.Note).Scan(&avis.ID)
	if err != nil {
		return avis, err
	}
	return avis, nil
}

// Update met à jour un avis client existant
func (r *AvisRepository) Update(avis model.AvisClient) error {
	// Vérifier si l'avis existe
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM avis_clients WHERE id = $1)", avis.ID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("avis non trouvé")
	}

	// Mettre à jour l'avis
	sqlStatement := `
	UPDATE avis_clients
	SET nom = $1, email = $2, avis = $3, note = $4
	WHERE id = $5`

	_, err = r.DB.Exec(sqlStatement, avis.Nom, avis.Email, avis.Avis, avis.Note, avis.ID)
	return err
}

// Delete supprime un avis client
func (r *AvisRepository) Delete(id int) error {
	// Vérifier si l'avis existe
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM avis_clients WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("avis non trouvé")
	}

	// Supprimer l'avis
	sqlStatement := "DELETE FROM avis_clients WHERE id = $1"
	_, err = r.DB.Exec(sqlStatement, id)
	return err
}