package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/polgarcia/golang-rest-api/internal/domain"
	"go.uber.org/zap"
)

// GetPokemonByName godoc
// @Summary Get Pokemon by name or ID
// @Description Get detailed information about a Pokemon by name or ID
// @Tags pokemon
// @Accept json
// @Produce json
// @Param nameOrId path string true "Pokemon name (e.g., 'pikachu') or ID (e.g., '25')"
// @Success 200 {object} domain.Pokemon
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 404 {object} ErrorResponse "Pokemon not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/pokemon/{nameOrId} [get]
func (h *Handler) GetPokemonByName(w http.ResponseWriter, r *http.Request) {
	nameOrID := chi.URLParam(r, "nameOrId")

	h.logger.Info("GetPokemonByName request",
		zap.String("name_or_id", nameOrID),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	// Get Pokemon from service
	pokemon, err := h.pokemonService.GetByName(r.Context(), nameOrID)
	if err != nil {
		h.handlePokemonError(w, err)
		return
	}

	// Return success response
	WriteJSON(w, http.StatusOK, pokemon, h.logger)
}

// ListPokemon godoc
// @Summary List Pokemon
// @Description Get a paginated list of Pokemon
// @Tags pokemon
// @Accept json
// @Produce json
// @Param limit query int false "Number of Pokemon to return (default: 20, max: 100)" default(20)
// @Param offset query int false "Number of Pokemon to skip (default: 0)" default(0)
// @Success 200 {object} domain.PokemonList
// @Failure 400 {object} ErrorResponse "Invalid query parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/pokemon [get]
func (h *Handler) ListPokemon(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Default values
	limit := 20
	offset := 0

	// Parse limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			h.logger.Debug("Invalid limit parameter", zap.String("limit", limitStr))
			WriteError(w, http.StatusBadRequest, "Invalid limit parameter", h.logger)
			return
		}
		limit = parsedLimit
	}

	// Parse offset
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			h.logger.Debug("Invalid offset parameter", zap.String("offset", offsetStr))
			WriteError(w, http.StatusBadRequest, "Invalid offset parameter", h.logger)
			return
		}
		offset = parsedOffset
	}

	h.logger.Info("ListPokemon request",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	// Get Pokemon list from service
	pokemonList, err := h.pokemonService.List(r.Context(), limit, offset)
	if err != nil {
		h.handlePokemonError(w, err)
		return
	}

	// Return success response
	WriteJSON(w, http.StatusOK, pokemonList, h.logger)
}

// handlePokemonError handles Pokemon-related errors and writes appropriate HTTP responses
func (h *Handler) handlePokemonError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrPokemonNotFound):
		WriteError(w, http.StatusNotFound, "Pokemon not found", h.logger)
	case errors.Is(err, domain.ErrInvalidInput):
		WriteError(w, http.StatusBadRequest, err.Error(), h.logger)
	case errors.Is(err, domain.ErrInvalidLimit):
		WriteError(w, http.StatusBadRequest, err.Error(), h.logger)
	case errors.Is(err, domain.ErrInvalidOffset):
		WriteError(w, http.StatusBadRequest, err.Error(), h.logger)
	case errors.Is(err, domain.ErrExternalAPI):
		h.logger.Error("External API error", zap.Error(err))
		WriteError(w, http.StatusBadGateway, "Failed to fetch data from external API", h.logger)
	default:
		h.logger.Error("Unexpected error", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Internal server error", h.logger)
	}
}
