package handler

import (
	"github.com/polgarcia/golang-rest-api/internal/domain"
	"github.com/polgarcia/golang-rest-api/pkg/logger"
)

// Handler is the base handler with dependencies
type Handler struct {
	pokemonService domain.PokemonService
	logger         *logger.Logger
}

// NewHandler creates a new handler with dependencies
func NewHandler(pokemonService domain.PokemonService, log *logger.Logger) *Handler {
	return &Handler{
		pokemonService: pokemonService,
		logger:         log,
	}
}
