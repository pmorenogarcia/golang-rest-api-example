package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/polgarcia/golang-rest-api/internal/domain"
	"github.com/polgarcia/golang-rest-api/pkg/logger"
	"go.uber.org/zap"
)

// PokemonService implements the domain.PokemonService interface
type PokemonService struct {
	client domain.PokemonClient
	logger *logger.Logger
}

// NewPokemonService creates a new Pokemon service
func NewPokemonService(client domain.PokemonClient, log *logger.Logger) *PokemonService {
	return &PokemonService{
		client: client,
		logger: log,
	}
}

// GetByName retrieves a Pokemon by name or ID
func (s *PokemonService) GetByName(ctx context.Context, nameOrID string) (*domain.Pokemon, error) {
	// Validate input
	if nameOrID == "" {
		s.logger.Debug("Invalid input: empty name or ID")
		return nil, fmt.Errorf("%w: name or ID cannot be empty", domain.ErrInvalidInput)
	}

	// Normalize name to lowercase
	nameOrID = strings.ToLower(strings.TrimSpace(nameOrID))

	s.logger.Info("Getting Pokemon",
		zap.String("name_or_id", nameOrID),
	)

	// Fetch Pokemon from client
	pokemon, err := s.client.FetchPokemon(ctx, nameOrID)
	if err != nil {
		s.logger.Error("Failed to get Pokemon",
			zap.String("name_or_id", nameOrID),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Info("Successfully retrieved Pokemon",
		zap.String("name", pokemon.Name),
		zap.Int("id", pokemon.ID),
	)

	return pokemon, nil
}
