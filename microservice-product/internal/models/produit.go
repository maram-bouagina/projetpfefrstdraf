package models

import (
	"time"

	"gorm.io/gorm"
)

type StatutProduit string
type VisibiliteProduit string

const (
	StatutBrouillon StatutProduit = "brouillon"
	StatutPublie    StatutProduit = "publie"
	StatutArchive   StatutProduit = "archive"

	VisibilitePublique VisibiliteProduit = "publique"
	VisibilitePrivee   VisibiliteProduit = "privee"
)

type Produit struct {
	ID              string            `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BoutiqueID      string            `gorm:"type:uuid;not null;index"                       json:"boutique_id"`
	Titre           string            `gorm:"type:varchar(255);not null"                     json:"titre"`
	Description     *string           `gorm:"type:text"                                      json:"description,omitempty"`
	Slug            string            `gorm:"type:varchar(255);not null;uniqueIndex:idx_slug_boutique" json:"slug"`
	Statut          StatutProduit     `gorm:"type:varchar(20);not null;default:brouillon"    json:"statut"`
	PrixDefaut      float64           `gorm:"type:decimal(12,4);not null;default:0"          json:"prix_defaut"`
	Devise          string            `gorm:"type:char(3);not null;default:EUR"              json:"devise"`
	SKU             *string           `gorm:"type:varchar(100)"                              json:"sku,omitempty"`
	SuiviStock      bool              `gorm:"not null;default:false"                         json:"suivi_stock"`
	QuantiteStock   int               `gorm:"not null;default:0"                             json:"quantite_stock"`
	Poids           *float64          `gorm:"type:decimal(10,4)"                             json:"poids,omitempty"`
	Dimensions      *string           `gorm:"type:varchar(100)"                              json:"dimensions,omitempty"`
	Marque          *string           `gorm:"type:varchar(255)"                              json:"marque,omitempty"`
	ClasseTaxe      *string           `gorm:"type:varchar(100)"                              json:"classe_taxe,omitempty"`
	Visibilite      VisibiliteProduit `gorm:"type:varchar(20);not null;default:publique"     json:"visibilite"`
	DatePublication *time.Time        `gorm:"type:timestamptz"                               json:"date_publication,omitempty"`
	SupprimeLe      gorm.DeletedAt    `gorm:"index"                                          json:"supprime_le,omitempty"`
	CreeLe          time.Time         `gorm:"autoCreateTime"                                 json:"cree_le"`
	MisAJourLe      time.Time         `gorm:"autoUpdateTime"                                 json:"mis_a_jour_le"`
}
