package model

type AvisClient struct {
	ID    int     `json:"id,omitempty"`
	Nom   string  `json:"nom"`
	Email string  `json:"email"`
	Avis  string  `json:"avis"`
	Note  float64 `json:"note"`
}

// ValidateAvis vérifie si l'avis est valide
func (a *AvisClient) ValidateAvis() bool {
	// Vérifie que les champs obligatoires ne sont pas vides
	if a.Nom == "" || a.Email == "" || a.Avis == "" {
		return false
	}

	// Vérifie que la note est dans l'intervalle valide
	if a.Note < 0 || a.Note > 5 {
		return false
	}

	return true
}