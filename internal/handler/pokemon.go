package handler

import (
	"errors"
	"net/http"

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

// GetPokemonCount godoc
// @Summary Get Pokemon count
// @Description Get the total number of Pokemon available
// @Tags pokemon
// @Accept json
// @Produce json
// @Success 200 {object} domain.PokemonCount
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/pokemon/count [get]
func (h *Handler) GetPokemonCount(w http.ResponseWriter, r *http.Request) {
	// Missing request logging (Claude should catch this)

	count, err := h.pokemonService.GetCount(r.Context())
	if err != nil {
		h.handlePokemonError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, count, h.logger)
}

// handlePokemonError handles Pokemon-related errors and writes appropriate HTTP responses
func (h *Handler) handlePokemonError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrPokemonNotFound):
		WriteError(w, http.StatusNotFound, "Pokemon not found", h.logger)
	case errors.Is(err, domain.ErrInvalidInput):
		WriteError(w, http.StatusBadRequest, err.Error(), h.logger)
	case errors.Is(err, domain.ErrExternalAPI):
		h.logger.Error("External API error", zap.Error(err))
		WriteError(w, http.StatusBadGateway, "Failed to fetch data from external API", h.logger)
	default:
		h.logger.Error("Unexpected error", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Internal server error", h.logger)
	}
}
