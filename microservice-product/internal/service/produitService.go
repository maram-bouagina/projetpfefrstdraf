package service

import (
	"context"
	"errors"
	"projet/internal/dto"
	"projet/internal/models"
	"projet/internal/repository"
	"strings"
	"time"
)

/*kik 3ada loula injectionde dépendance ou thneya constructeur*/
type ProduitService struct {
	repo *repository.ProduitRepo
}

func NewProduitService(repo *repository.ProduitRepo) *ProduitService {
	return &ProduitService{repo: repo}
}

/*t7wl ll produit model il reponse illi fo dto*/
func (s *ProduitService) toResponse(p models.Produit) dto.ProduitResponse {
	// 1. Convertir les options en réponses
	options := make([]dto.OptionProduitResponse, len(p.Options))
	for i, opt := range p.Options {
		// Convertir les valeurs de cette option
		valeurs := make([]dto.ValeurOptionResponse, len(opt.ValeurOpts))
		for j, val := range opt.ValeurOpts {
			valeurs[j] = dto.ValeurOptionResponse{
				ID:       val.ID,
				OptionID: val.OptionID,
				Valeur:   val.Valeur,
				Position: val.Position,
			}
		}

		options[i] = dto.OptionProduitResponse{
			ID:         opt.ID,
			ProduitID:  opt.ProduitID,
			Nom:        opt.Nom,
			Position:   opt.Position,
			CreeLe:     opt.CreeLe,
			MisAJourLe: opt.MisAJourLe,
			ValeurOpts: valeurs,
		}
	}

	return dto.ProduitResponse{
		ID:              p.ID,
		Titre:           p.Titre,
		Description:     p.Description,
		Slug:            p.Slug,
		Statut:          p.Statut,
		PrixDefaut:      p.PrixDefaut,
		Devise:          p.Devise,
		SKU:             p.SKU,
		SuiviStock:      p.SuiviStock,
		QuantiteStock:   p.QuantiteStock,
		Poids:           p.Poids,
		Dimensions:      p.Dimensions,
		Marque:          p.Marque,
		ClasseTaxe:      p.ClasseTaxe,
		Visibilite:      p.Visibilite,
		DatePublication: p.DatePublication,
		CreeLe:          p.CreeLe,
		MisAJourLe:      p.MisAJourLe,
		Options:         options,
		Variantes:       []dto.VarianteResponse{},
	}
}

func (s *ProduitService) Create(ctx context.Context, boutiqueID string, req dto.RequeteCreationProduit) (*dto.ProduitResponse, error) {
	if boutiqueID == "" {
		return nil, errors.New("boutique ID is required")
	}

	//idhe knou much mwjoud ou idha knu fergh
	if req.Slug == nil || *req.Slug == "" {

		//ticrétih bmt3 lvaleur mtaa titre
		slug := strings.ToLower(strings.ReplaceAll(req.Titre, " ", "-"))
		req.Slug = &slug
	}

	//tisnaa3 produit
	produit := &models.Produit{
		BoutiqueID:      boutiqueID,
		Titre:           req.Titre,
		Description:     req.Description,
		Slug:            *req.Slug,
		Statut:          req.Statut,
		PrixDefaut:      req.PrixDefaut,
		Devise:          req.Devise,
		SKU:             req.SKU,
		SuiviStock:      req.SuiviStock,
		QuantiteStock:   req.QuantiteStock,
		Poids:           req.Poids,
		Dimensions:      req.Dimensions,
		Marque:          req.Marque,
		ClasseTaxe:      req.ClasseTaxe,
		Visibilite:      req.Visibilite,
		DatePublication: req.DatePublication,
	}

	//t3yt li repositroy
	created, err := s.repo.CreateProduit(ctx, produit)
	if err != nil {
		return nil, err
	}
	//t3yt ll helper (func tit3wd bech nhiw redendance) illi lfou9
	resp := s.toResponse(*created)
	return &resp, nil
}

func (s *ProduitService) List(ctx context.Context, boutiqueID string) ([]dto.ProduitResponse, error) {
	if boutiqueID == "" {
		return nil, errors.New("boutique ID is required")
	}
	//t3yt li repo
	produits, err := s.repo.ListProduits(ctx, boutiqueID)
	if err != nil {
		return nil, err
	}

	//bch naamlou slice : tableaux dynamiques mi liste illi khdhineha
	//resp = [ {}, {}, {} ]
	//taille mta3ha nb de produits dans le db
	resp := make([]dto.ProduitResponse, len(produits))
	//l'index ou kl produit
	for i, p := range produits {
		resp[i] = s.toResponse(p)
	}
	return resp, nil
}

func (s *ProduitService) GetByID(ctx context.Context, id, boutiqueID string) (*dto.ProduitResponse, error) {
	if boutiqueID == "" {
		return nil, errors.New("boutique ID is required")
	}
	produit, err := s.repo.GetByID(ctx, id, boutiqueID)
	if err != nil {
		return nil, err
	}
	if produit == nil {
		return nil, errors.New("product not found")
	}
	resp := s.toResponse(*produit)
	return &resp, nil
}

func (s *ProduitService) Update(ctx context.Context, id, boutiqueID string, req dto.RequeteUpdateProduit) (*dto.ProduitResponse, error) {
	if boutiqueID == "" {
		return nil, errors.New("boutique ID is required")
	}

	//t3yt ll func illi fou9ha
	_, err := s.GetByID(ctx, id, boutiqueID)
	if err != nil {
		return nil, err
	}
	//tisnaa slice
	updates := make(map[string]interface{})

	//affectation des valeurs
	if req.Titre != nil {
		updates["titre"] = *req.Titre
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Statut != nil {
		updates["statut"] = *req.Statut
	}
	if req.PrixDefaut != nil {
		updates["prix_defaut"] = *req.PrixDefaut
	}
	if req.Devise != nil {
		updates["devise"] = *req.Devise
	}
	if req.SKU != nil {
		updates["sku"] = *req.SKU
	}
	if req.SuiviStock != nil {
		updates["suivi_stock"] = *req.SuiviStock
	}
	if req.QuantiteStock != nil {
		updates["quantite_stock"] = *req.QuantiteStock
	}
	if req.Poids != nil {
		updates["poids"] = *req.Poids
	}
	if req.Dimensions != nil {
		updates["dimensions"] = *req.Dimensions
	}
	if req.Marque != nil {
		updates["marque"] = *req.Marque
	}
	if req.ClasseTaxe != nil {
		updates["classe_taxe"] = *req.ClasseTaxe
	}
	if req.Visibilite != nil {
		updates["visibilite"] = *req.Visibilite
	}
	if req.DatePublication != nil {
		updates["date_publication"] = *req.DatePublication
	}
	updates["mis_a_jour_le"] = time.Now()

	updated, err := s.repo.Update(ctx, id, boutiqueID, updates)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, errors.New("product not found after update")
	}
	resp := s.toResponse(*updated)
	return &resp, nil
}

func (s *ProduitService) Delete(ctx context.Context, id, boutiqueID string) error {
	if boutiqueID == "" {
		return errors.New("boutique ID is required")
	}
	deleted, err := s.repo.DeleteById(ctx, id, boutiqueID)
	if err != nil {
		return err
	}
	if !deleted {
		return errors.New("product not found")
	}
	return nil
}

func (s *ProduitService) Search(ctx context.Context, boutiqueID string, filter dto.FiltreProduit) ([]dto.ProduitResponse, int, int, error) {
	if boutiqueID == "" {
		return nil, 0, 0, errors.New("boutique ID is required")
	}

	/*fi go ki naamlouch valeur l valeur par défaut mtaa les entier est 0*/

	if filter.Page == 0 {
		filter.Page = 1
	}
	//ki mayaatikch twalli 20 par defaut
	if filter.Limite == 0 {
		filter.Limite = 20
	}

	produits, err := s.repo.GetWithFilter(ctx, boutiqueID, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	resp := make([]dto.ProduitResponse, len(produits))
	for i, p := range produits {
		resp[i] = s.toResponse(p)
	}
	return resp, filter.Page, filter.Limite, nil
}
