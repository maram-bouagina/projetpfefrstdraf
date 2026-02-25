package models

import "time"

type Variante struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProduitID     string    `gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE;references:produits(id)" json:"produit_id"`
	SKU           string    `gorm:"type:varchar(100);not null;uniqueIndex"         json:"sku"`
	Prix          *float64  `gorm:"type:decimal(12,4)"                             json:"prix,omitempty"`
	QuantiteStock int       `gorm:"not null;default:0"                             json:"quantite_stock"`
	CodeBarres    *string   `gorm:"type:varchar(100)"                              json:"code_barres,omitempty"`
	Poids         *float64  `gorm:"type:decimal(10,4)"                             json:"poids,omitempty"`
	Images        []string  `gorm:"type:text[];serializer:json"                    json:"images,omitempty"`
	CreeLe        time.Time `gorm:"autoCreateTime"                                 json:"cree_le"`
	MisAJourLe    time.Time `gorm:"autoUpdateTime"                                 json:"mis_a_jour_le"`

	// Relations
	ValeurOptions []ValeurOption `gorm:"many2many:variante_valeur_option;" json:"valeur_options,omitempty"`
}

type VarianteValeurOption struct {
	VarianteID     string `gorm:"type:uuid;primaryKey" json:"variante_id"`
	ValeurOptionID string `gorm:"type:uuid;primaryKey" json:"valeur_option_id"`
}
