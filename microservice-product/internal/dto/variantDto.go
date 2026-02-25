package dto

import "time"

type RequeteCreationVariante struct {
	SKU           string   `json:"sku"             validate:"required,min=1,max=100"`
	Prix          *float64 `json:"prix"            validate:"omitempty,min=0"`
	QuantiteStock int      `json:"quantite_stock"  validate:"min=0"`
	CodeBarres    *string  `json:"code_barres"`
	Poids         *float64 `json:"poids"           validate:"omitempty,min=0"`
	Images        []string `json:"images"`

	ValeurOptionIDs []string `json:"valeur_option_ids" validate:"required,min=1"`
}

type RequeteUpdateVariante struct {
	SKU             *string  `json:"sku"             validate:"omitempty,min=1,max=100"`
	Prix            *float64 `json:"prix"            validate:"omitempty,min=0"`
	QuantiteStock   *int     `json:"quantite_stock"  validate:"omitempty,min=0"`
	CodeBarres      *string  `json:"code_barres"`
	Poids           *float64 `json:"poids"           validate:"omitempty,min=0"`
	Images          []string `json:"images"`
	ValeurOptionIDs []string `json:"valeur_option_ids" validate:"omitempty,min=1"`
}

type VarianteResponse struct {
	ID            string                 `json:"id"`
	ProduitID     string                 `json:"produit_id"`
	SKU           string                 `json:"sku"`
	Prix          *float64               `json:"prix,omitempty"`
	QuantiteStock int                    `json:"quantite_stock"`
	CodeBarres    *string                `json:"code_barres,omitempty"`
	Poids         *float64               `json:"poids,omitempty"`
	Images        []string               `json:"images,omitempty"`
	CreeLe        time.Time              `json:"cree_le"`
	MisAJourLe    time.Time              `json:"mis_a_jour_le"`
	ValeurOptions []ValeurOptionResponse `json:"valeur_options,omitempty"`
	// Prix effectif = Prix si présent, sinon PrixDefaut du produit
	PrixEffectif float64 `json:"prix_effectif"`
}
