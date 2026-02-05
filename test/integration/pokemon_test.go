package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/polgarcia/golang-rest-api/internal/client"
	"github.com/polgarcia/golang-rest-api/internal/domain"
	"github.com/polgarcia/golang-rest-api/internal/handler"
	"github.com/polgarcia/golang-rest-api/internal/server"
	"github.com/polgarcia/golang-rest-api/internal/service"
	"github.com/polgarcia/golang-rest-api/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestServer creates a test server with all dependencies
func setupTestServer(t *testing.T) http.Handler {
	// Create test logger
	log, err := logger.New("error", "console")
	require.NoError(t, err)

	// Create PokeAPI client
	pokemonClient := client.NewPokeAPIClient(
		"https://pokeapi.co/api/v2",
		30000000000, // 30 seconds
		log,
	)

	// Create Pokemon service
	pokemonService := service.NewPokemonService(pokemonClient, log)

	// Create handlers
	h := handler.NewHandler(pokemonService, log)

	// Setup routes
	router := server.SetupRoutes(h, log, "*")

	return router
}

func TestHealthCheckEndpoint(t *testing.T) {
	router := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response handler.HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "ok", response.Status)
	assert.NotZero(t, response.Timestamp)
}

func TestGetPokemonByName(t *testing.T) {
	router := setupTestServer(t)

	tests := []struct {
		name           string
		pokemonName    string
		expectedStatus int
		checkResponse  func(t *testing.T, pokemon *domain.Pokemon)
	}{
		{
			name:           "Get Pikachu",
			pokemonName:    "pikachu",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, pokemon *domain.Pokemon) {
				assert.Equal(t, 25, pokemon.ID)
				assert.Equal(t, "pikachu", pokemon.Name)
				assert.NotEmpty(t, pokemon.Types)
				assert.NotEmpty(t, pokemon.Abilities)
			},
		},
		{
			name:           "Get Pokemon by ID",
			pokemonName:    "1",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, pokemon *domain.Pokemon) {
				assert.Equal(t, 1, pokemon.ID)
				assert.Equal(t, "bulbasaur", pokemon.Name)
			},
		},
		{
			name:           "Pokemon not found",
			pokemonName:    "nonexistent123",
			expectedStatus: http.StatusNotFound,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/pokemon/"+tt.pokemonName, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil && w.Code == http.StatusOK {
				var pokemon domain.Pokemon
				err := json.NewDecoder(w.Body).Decode(&pokemon)
				require.NoError(t, err)

				tt.checkResponse(t, &pokemon)
			}
		})
	}
}

func TestCORSHeaders(t *testing.T) {
	router := setupTestServer(t)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/pokemon", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestRequestIDHeader(t *testing.T) {
	router := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	requestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID, "X-Request-ID header should be present")
}

// Benchmark tests
func BenchmarkGetPokemon(b *testing.B) {
	log, _ := logger.New("error", "console")
	pokemonClient := client.NewPokeAPIClient("https://pokeapi.co/api/v2", 30000000000, log)
	pokemonService := service.NewPokemonService(pokemonClient, log)
	h := handler.NewHandler(pokemonService, log)
	router := server.SetupRoutes(h, log, "*")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/pokemon/pikachu", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

