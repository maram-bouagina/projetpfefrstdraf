package models

import "time"

type OptionProduit struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProduitID  string    `gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE;references:produits(id)" json:"produit_id"`
	Nom        string    `gorm:"type:varchar(100);not null"                     json:"nom"`
	Position   int       `gorm:"not null;default:0"                             json:"position"`
	CreeLe     time.Time `gorm:"autoCreateTime"                                 json:"cree_le"`
	MisAJourLe time.Time `gorm:"autoUpdateTime"                                 json:"mis_a_jour_le"`

	// Relations
	ValeurOpts []ValeurOption `gorm:"foreignKey:OptionID;constraint:OnDelete:CASCADE" json:"valeur_opts,omitempty"`
}

type ValeurOption struct {
	ID       string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OptionID string `gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE;references:option_produits(id)" json:"option_id"`
	Valeur   string `gorm:"type:varchar(100);not null"                     json:"valeur"`
	Position int    `gorm:"not null;default:0"                             json:"position"`
}
