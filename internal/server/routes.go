package server

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/polgarcia/golang-rest-api/internal/handler"
	"github.com/polgarcia/golang-rest-api/internal/middleware"
	"github.com/polgarcia/golang-rest-api/pkg/logger"
)

// SetupRoutes configures all application routes
func SetupRoutes(h *handler.Handler, log *logger.Logger, corsOrigins string) *chi.Mux {
	r := chi.NewRouter()

	// Apply middleware chain
	r.Use(middleware.Recovery(log))
	r.Use(middleware.Logger(log))
	r.Use(middleware.CORS(corsOrigins))

	// Health check endpoint (no prefix)
	r.Get("/health", h.HealthCheck)

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Pokemon endpoints
		r.Route("/pokemon", func(r chi.Router) {
			r.Get("/{nameOrId}", h.GetPokemonByName)
		})
	})

	return r
}
