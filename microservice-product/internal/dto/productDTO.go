package dto

import (
	"projet/internal/models"
	"time"
)

type RequeteCreationProduit struct {
	Titre           string                   `json:"titre"           validate:"required,min=1,max=255"`
	Description     *string                  `json:"description"`
	Slug            *string                  `json:"slug"`
	Statut          models.StatutProduit     `json:"statut"          validate:"required,oneof=brouillon publie archive"`
	PrixDefaut      float64                  `json:"prix_defaut"     validate:"min=0"`
	Devise          string                   `json:"devise"          validate:"required,len=3"`
	SKU             *string                  `json:"sku"`
	SuiviStock      bool                     `json:"suivi_stock"`
	QuantiteStock   int                      `json:"quantite_stock"  validate:"min=0"`
	Poids           *float64                 `json:"poids"           validate:"omitempty,min=0"`
	Dimensions      *string                  `json:"dimensions"`
	Marque          *string                  `json:"marque"`
	ClasseTaxe      *string                  `json:"classe_taxe"`
	Visibilite      models.VisibiliteProduit `json:"visibilite"      validate:"required,oneof=publique privee"`
	DatePublication *time.Time               `json:"date_publication"`
}

type RequeteUpdateProduit struct {
	Titre           *string                   `json:"titre"           validate:"omitempty,min=1,max=255"`
	Description     *string                   `json:"description"`
	Slug            *string                   `json:"slug"`
	Statut          *models.StatutProduit     `json:"statut"          validate:"omitempty,oneof=brouillon publie archive"`
	PrixDefaut      *float64                  `json:"prix_defaut"     validate:"omitempty,min=0"`
	Devise          *string                   `json:"devise"          validate:"omitempty,len=3"`
	SKU             *string                   `json:"sku"`
	SuiviStock      *bool                     `json:"suivi_stock"`
	QuantiteStock   *int                      `json:"quantite_stock"  validate:"omitempty,min=0"`
	Poids           *float64                  `json:"poids"           validate:"omitempty,min=0"`
	Dimensions      *string                   `json:"dimensions"`
	Marque          *string                   `json:"marque"`
	ClasseTaxe      *string                   `json:"classe_taxe"`
	Visibilite      *models.VisibiliteProduit `json:"visibilite"      validate:"omitempty,oneof=publique privee"`
	DatePublication *time.Time                `json:"date_publication"`
}

type ProduitResponse struct {
	ID              string                   `json:"id"`
	Titre           string                   `json:"titre"`
	Description     *string                  `json:"description,omitempty"`
	Slug            string                   `json:"slug"`
	Statut          models.StatutProduit     `json:"statut"`
	PrixDefaut      float64                  `json:"prix_defaut"`
	Devise          string                   `json:"devise"`
	SKU             *string                  `json:"sku,omitempty"`
	SuiviStock      bool                     `json:"suivi_stock"`
	QuantiteStock   int                      `json:"quantite_stock"`
	Poids           *float64                 `json:"poids,omitempty"`
	Dimensions      *string                  `json:"dimensions,omitempty"`
	Marque          *string                  `json:"marque,omitempty"`
	ClasseTaxe      *string                  `json:"classe_taxe,omitempty"`
	Visibilite      models.VisibiliteProduit `json:"visibilite"`
	DatePublication *time.Time               `json:"date_publication,omitempty"`
	CreeLe          time.Time                `json:"created_at"`
	MisAJourLe      time.Time                `json:"updated_at"`
}

type FiltreProduit struct {
	Statut          *models.StatutProduit     `query:"statut"`
	Visibilite      *models.VisibiliteProduit `query:"visibilite"`
	Marque          *string                   `query:"marque"`
	Recherche       *string                   `query:"search"`
	Page            int                       `query:"page"`
	Limite          int                       `query:"limit"`
	InclureSupprime bool                      `query:"inclure_supprime"`
}
