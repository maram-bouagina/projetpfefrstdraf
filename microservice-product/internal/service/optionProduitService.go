package service

import (
	"context"
	"errors"
	"fmt"
	"projet/internal/dto"
	"projet/internal/models"
	"projet/internal/repository"
	"time"
)

type OptionProduitService struct {
	repo *repository.OptionProduitValeurRepo
}

func NewOptionProduitService(repo *repository.OptionProduitValeurRepo) *OptionProduitService {
	return &OptionProduitService{repo: repo}
}

// ------------------------------------------------------------
// Convertisseurs
// ------------------------------------------------------------
func (s *OptionProduitService) toResponseValeurOpt(vo models.ValeurOption) dto.ValeurOptionResponse {
	return dto.ValeurOptionResponse{
		ID:       vo.ID,
		OptionID: vo.OptionID,
		Valeur:   vo.Valeur,
		Position: vo.Position,
	}
}

func (s *OptionProduitService) toResponseOptProd(po models.OptionProduit) dto.OptionProduitResponse {
	vopts := make([]dto.ValeurOptionResponse, len(po.ValeurOpts))
	for i, vopt := range po.ValeurOpts {
		vopts[i] = s.toResponseValeurOpt(vopt)
	}

	return dto.OptionProduitResponse{
		ID:         po.ID,
		ProduitID:  po.ProduitID,
		Nom:        po.Nom,
		Position:   po.Position,
		CreeLe:     po.CreeLe,
		MisAJourLe: po.MisAJourLe,
		ValeurOpts: vopts,
	}
}

// ------------------------------------------------------------
// Créer une option avec ses valeurs (les valeurs sont créées séparément)
// ------------------------------------------------------------
func (s *OptionProduitService) CreationOptionProduit(
	ctx context.Context,
	produitID string,
	req dto.RequeteCreationOption,
) (*dto.OptionProduitResponse, error) {

	// Gérer la position
	position := req.Position
	if position == 0 {
		count, err := s.repo.CountOptionsByProduit(ctx, produitID)
		if err != nil {
			return nil, fmt.Errorf("impossible de compter les options: %v", err)
		}
		position = count + 1
	}

	// Créer l'option seulement (sans valeurs)
	nouvelleOption := &models.OptionProduit{
		ProduitID:  produitID,
		Nom:        req.Nom,
		Position:   position,
		CreeLe:     time.Now(),
		MisAJourLe: time.Now(),
	}

	cree, err := s.repo.CreationOptProduit(ctx, nouvelleOption)
	if err != nil {
		return nil, fmt.Errorf("échec de la création: %v", err)
	}

	// Retourner l'option (sans valeurs pour l'instant)
	reponse := s.toResponseOptProd(*cree)
	return &reponse, nil
}

// ------------------------------------------------------------
// Créer une valeur d'option pour une option existante
// ------------------------------------------------------------
func (s *OptionProduitService) CreationValeurOption(
	ctx context.Context,
	optionID string,
	req dto.RequeteCreationValeurOption,
) (*dto.ValeurOptionResponse, error) {

	// Vérifier que l'option existe
	option, err := s.repo.GetByIdOptionproduit(ctx, optionID)
	if err != nil {
		return nil, err
	}
	if option == nil {
		return nil, errors.New("option non trouvée")
	}

	// Gérer la position
	position := req.Position
	if position == 0 {
		count, err := s.repo.CountValeursByOption(ctx, optionID)
		if err != nil {
			return nil, fmt.Errorf("impossible de compter les valeurs: %v", err)
		}
		position = count + 1
	}

	// Créer la valeur
	nouvelleValeur := &models.ValeurOption{
		OptionID: optionID,
		Valeur:   req.Valeur,
		Position: position,
	}

	cree, err := s.repo.CreationValeurOption(ctx, nouvelleValeur)
	if err != nil {
		return nil, fmt.Errorf("échec de la création de la valeur: %v", err)
	}

	reponse := s.toResponseValeurOpt(*cree)
	return &reponse, nil
}

// ------------------------------------------------------------
// Lister toutes les options d'un produit (avec leurs valeurs)
// ------------------------------------------------------------
func (s *OptionProduitService) ListOptionProduit(ctx context.Context, produitID string) ([]dto.OptionProduitResponse, error) {
	if produitID == "" {
		return nil, errors.New("ID du produit requis")
	}

	options, err := s.repo.ListeOptProduits(ctx, produitID)
	if err != nil {
		return nil, err
	}

	// Pour chaque option, on va chercher ses valeurs
	resultats := make([]dto.OptionProduitResponse, len(options))
	for i, opt := range options {
		// Récupérer les valeurs de cette option
		valeurs, err := s.repo.ListeValeursOption(ctx, opt.ID)
		if err != nil {
			return nil, err
		}
		opt.ValeurOpts = valeurs
		resultats[i] = s.toResponseOptProd(opt)
	}

	return resultats, nil
}

// ------------------------------------------------------------
// Récupérer une option spécifique par son ID (avec ses valeurs)
// ------------------------------------------------------------
func (s *OptionProduitService) GetByIDOptionProduit(ctx context.Context, id string) (*dto.OptionProduitResponse, error) {
	option, err := s.repo.GetByIdOptionproduit(ctx, id)
	if err != nil {
		return nil, err
	}
	if option == nil {
		return nil, errors.New("option non trouvée")
	}

	// Récupérer les valeurs de cette option
	valeurs, err := s.repo.ListeValeursOption(ctx, option.ID)
	if err != nil {
		return nil, err
	}
	option.ValeurOpts = valeurs

	reponse := s.toResponseOptProd(*option)
	return &reponse, nil
}

// ------------------------------------------------------------
// Mettre à jour une option (nom, position)
// ------------------------------------------------------------
func (s *OptionProduitService) Update(ctx context.Context, id string, req dto.RequeteUpdateOption) (*dto.OptionProduitResponse, error) {
	// Vérifier que l'option existe
	if _, err := s.GetByIDOptionProduit(ctx, id); err != nil {
		return nil, err
	}

	modifications := make(map[string]interface{})
	if req.Nom != nil {
		modifications["nom"] = *req.Nom
	}
	if req.Position != nil {
		modifications["position"] = *req.Position
	}
	modifications["mis_a_jour_le"] = time.Now()

	if len(modifications) == 0 {
		return s.GetByIDOptionProduit(ctx, id)
	}

	_, err := s.repo.UpdateProduitOption(ctx, id, modifications)
	if err != nil {
		return nil, err
	}

	return s.GetByIDOptionProduit(ctx, id)
}

// ------------------------------------------------------------
// Mettre à jour une valeur d'option
// ------------------------------------------------------------
func (s *OptionProduitService) UpdateValeur(
	ctx context.Context,
	id string,
	req dto.RequeteUpdateValeurOption,
) (*dto.ValeurOptionResponse, error) {

	// Vérifier que la valeur existe
	valeur, err := s.repo.GetByIDValeurOption(ctx, id)
	if err != nil {
		return nil, err
	}
	if valeur == nil {
		return nil, errors.New("valeur non trouvée")
	}

	modifications := make(map[string]interface{})
	if req.Valeur != nil {
		modifications["valeur"] = *req.Valeur
	}
	if req.Position != nil {
		modifications["position"] = *req.Position
	}

	if len(modifications) == 0 {
		reponse := s.toResponseValeurOpt(*valeur)
		return &reponse, nil
	}

	updated, err := s.repo.UpdateValeurOpt(ctx, id, modifications)
	if err != nil {
		return nil, err
	}

	reponse := s.toResponseValeurOpt(*updated)
	return &reponse, nil
}

// ------------------------------------------------------------
// Supprimer une option (et ses valeurs par CASCADE)
// ------------------------------------------------------------
func (s *OptionProduitService) Delete(ctx context.Context, id string) error {
	supprime, err := s.repo.SupprimerOptPById(ctx, id)
	if err != nil {
		return err
	}
	if !supprime {
		return errors.New("option non trouvée")
	}
	return nil
}

// ------------------------------------------------------------
// Supprimer une valeur d'option
// ------------------------------------------------------------
func (s *OptionProduitService) DeleteValeur(ctx context.Context, id string) error {
	supprime, err := s.repo.SupprimerByIdVOpt(ctx, id)
	if err != nil {
		return err
	}
	if !supprime {
		return errors.New("valeur non trouvée")
	}
	return nil
}

// ------------------------------------------------------------
// Lister les valeurs d'une option
// ------------------------------------------------------------
func (s *OptionProduitService) ListValeursByOption(ctx context.Context, optionID string) ([]dto.ValeurOptionResponse, error) {
	valeurs, err := s.repo.ListeValeursOption(ctx, optionID)
	if err != nil {
		return nil, err
	}

	resultats := make([]dto.ValeurOptionResponse, len(valeurs))
	for i, v := range valeurs {
		resultats[i] = s.toResponseValeurOpt(v)
	}
	return resultats, nil
}
