package dto

import "time"

type RequeteCreationOption struct {
	Nom      string `json:"nom"      validate:"required,min=1,max=100"`
	Position int    `json:"position" validate:"min=0"`
}

type RequeteUpdateOption struct {
	Nom      *string `json:"nom"      validate:"omitempty,min=1,max=100"`
	Position *int    `json:"position" validate:"omitempty,min=0"`
}

type OptionProduitResponse struct {
	ID         string                 `json:"id"`
	ProduitID  string                 `json:"produit_id"`
	Nom        string                 `json:"nom"`
	Position   int                    `json:"position"`
	CreeLe     time.Time              `json:"cree_le"`
	MisAJourLe time.Time              `json:"mis_a_jour_le"`
	ValeurOpts []ValeurOptionResponse `json:"valeur_opts,omitempty"`
}

type RequeteCreationValeurOption struct {
	Valeur   string `json:"valeur"    validate:"required,min=1,max=100"`
	Position int    `json:"position"  validate:"min=0"`
}

type RequeteUpdateValeurOption struct {
	Valeur   *string `json:"valeur"    validate:"omitempty,min=1,max=100"`
	Position *int    `json:"position"  validate:"omitempty,min=0"`
}

type ValeurOptionResponse struct {
	ID       string `json:"id"`
	OptionID string `json:"option_id"`
	Valeur   string `json:"valeur"`
	Position int    `json:"position"`
}
