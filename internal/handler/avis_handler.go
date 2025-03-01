package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"api-avis-clients/internal/model"
	"api-avis-clients/internal/repository"
)

// AvisHandler gère les requêtes HTTP liées aux avis clients
type AvisHandler struct {
	repo *repository.AvisRepository
}

// NewAvisHandler crée une nouvelle instance de AvisHandler
func NewAvisHandler(repo *repository.AvisRepository) *AvisHandler {
	return &AvisHandler{repo: repo}
}

// GetAll renvoie tous les avis clients
func (h *AvisHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	avisList, err := h.repo.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, avisList)
}

// CreateMultiple ajoute plusieurs avis clients
func (h *AvisHandler) CreateMultiple(w http.ResponseWriter, r *http.Request) {
	var avisList []model.AvisClient
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&avisList); err != nil {
		respondWithError(w, http.StatusBadRequest, "Requête invalide")
		return
	}
	defer r.Body.Close()

	// Validation des avis
	for _, avis := range avisList {
		if !avis.ValidateAvis() {
			respondWithError(w, http.StatusBadRequest, "Données d'avis invalides")
			return
		}
	}

	// Insérer tous les avis
	var createdAvisList []model.AvisClient
	for _, avis := range avisList {
		createdAvis, err := h.repo.Create(avis)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		createdAvisList = append(createdAvisList, createdAvis)
	}

	respondWithJSON(w, http.StatusCreated, createdAvisList)
}


func (h *AvisHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID invalide")
		return
	}

	avis, err := h.repo.GetByID(id)
	if err != nil {
		if err.Error() == "avis non trouvé" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, avis)
}

// Create ajoute un nouvel avis client
func (h *AvisHandler) Create(w http.ResponseWriter, r *http.Request) {
	var avis model.AvisClient
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&avis); err != nil {
		respondWithError(w, http.StatusBadRequest, "Requête invalide")
		return
	}
	defer r.Body.Close()

	// Validation
	if !avis.ValidateAvis() {
		respondWithError(w, http.StatusBadRequest, "Données d'avis invalides")
		return
	}

	createdAvis, err := h.repo.Create(avis)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, createdAvis)
}

// Update met à jour un avis client existant
func (h *AvisHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID invalide")
		return
	}

	var avis model.AvisClient
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&avis); err != nil {
		respondWithError(w, http.StatusBadRequest, "Requête invalide")
		return
	}
	defer r.Body.Close()

	avis.ID = id

	// Validation
	if !avis.ValidateAvis() {
		respondWithError(w, http.StatusBadRequest, "Données d'avis invalides")
		return
	}

	err = h.repo.Update(avis)
	if err != nil {
		if err.Error() == "avis non trouvé" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, avis)
}

// Delete supprime un avis client
func (h *AvisHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID invalide")
		return
	}

	err = h.repo.Delete(id)
	if err != nil {
		if err.Error() == "avis non trouvé" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Avis supprimé avec succès"})
}

// Fonctions utilitaires pour les réponses JSON
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erreur lors de la sérialisation JSON"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}