package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/polgarcia/golang-rest-api/internal/domain"
	"github.com/polgarcia/golang-rest-api/pkg/logger"
	"go.uber.org/zap"
)

const (
	// MaxLimit is the maximum number of Pokemon that can be requested at once
	MaxLimit = 100
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

// List retrieves a paginated list of Pokemon
func (s *PokemonService) List(ctx context.Context, limit, offset int) (*domain.PokemonList, error) {
	// Validate limit
	if limit <= 0 {
		s.logger.Debug("Invalid limit", zap.Int("limit", limit))
		return nil, domain.ErrInvalidLimit
	}

	if limit > MaxLimit {
		s.logger.Debug("Limit exceeds maximum",
			zap.Int("limit", limit),
			zap.Int("max_limit", MaxLimit),
		)
		return nil, domain.ErrInvalidLimit
	}

	// Validate offset
	if offset < 0 {
		s.logger.Debug("Invalid offset", zap.Int("offset", offset))
		return nil, domain.ErrInvalidOffset
	}

	s.logger.Info("Listing Pokemon",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	// Fetch Pokemon list from client
	pokemonList, err := s.client.FetchPokemonList(ctx, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list Pokemon",
			zap.Int("limit", limit),
			zap.Int("offset", offset),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Info("Successfully retrieved Pokemon list",
		zap.Int("count", pokemonList.Count),
		zap.Int("results", len(pokemonList.Results)),
	)

	return pokemonList, nil
}
