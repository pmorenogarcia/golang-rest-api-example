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

func TestListPokemon(t *testing.T) {
	router := setupTestServer(t)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkResponse  func(t *testing.T, list *domain.PokemonList)
	}{
		{
			name:           "List default (no params)",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, list *domain.PokemonList) {
				assert.Greater(t, list.Count, 0, "Count should be greater than 0")
				assert.NotEmpty(t, list.Results, "Results should not be empty")
				assert.LessOrEqual(t, len(list.Results), 20, "Default limit should be 20")
				// Verify first result has expected fields
				if len(list.Results) > 0 {
					assert.NotEmpty(t, list.Results[0].Name)
					assert.NotEmpty(t, list.Results[0].URL)
				}
			},
		},
		{
			name:           "List with custom limit",
			queryParams:    "?limit=5&offset=0",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, list *domain.PokemonList) {
				assert.Len(t, list.Results, 5, "Should return exactly 5 results")
				assert.Greater(t, list.Count, 0)
			},
		},
		{
			name:           "List with pagination",
			queryParams:    "?limit=10&offset=10",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, list *domain.PokemonList) {
				assert.LessOrEqual(t, len(list.Results), 10, "Should return at most 10 results")
				assert.Greater(t, list.Count, 0)
				// Verify pagination - results should be different from offset 0
				if len(list.Results) > 0 {
					assert.NotEqual(t, "bulbasaur", list.Results[0].Name, "First result should not be bulbasaur (offset 10)")
				}
			},
		},
		{
			name:           "List first page with limit 3",
			queryParams:    "?limit=3&offset=0",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, list *domain.PokemonList) {
				assert.Len(t, list.Results, 3)
				// First three should be bulbasaur, ivysaur, venusaur
				assert.Equal(t, "bulbasaur", list.Results[0].Name)
			},
		},
		{
			name:           "Invalid limit (too high)",
			queryParams:    "?limit=200",
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name:           "Invalid limit (zero)",
			queryParams:    "?limit=0",
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name:           "Invalid limit (negative)",
			queryParams:    "?limit=-1",
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name:           "Invalid offset (negative)",
			queryParams:    "?offset=-1",
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name:           "Invalid limit (not a number)",
			queryParams:    "?limit=abc",
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name:           "Invalid offset (not a number)",
			queryParams:    "?offset=xyz",
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name:           "Max valid limit (100)",
			queryParams:    "?limit=100&offset=0",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, list *domain.PokemonList) {
				assert.LessOrEqual(t, len(list.Results), 100)
				assert.Greater(t, list.Count, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/pokemon"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil && w.Code == http.StatusOK {
				var list domain.PokemonList
				err := json.NewDecoder(w.Body).Decode(&list)
				require.NoError(t, err)

				tt.checkResponse(t, &list)
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

func BenchmarkListPokemon(b *testing.B) {
	log, _ := logger.New("error", "console")
	pokemonClient := client.NewPokeAPIClient("https://pokeapi.co/api/v2", 30000000000, log)
	pokemonService := service.NewPokemonService(pokemonClient, log)
	h := handler.NewHandler(pokemonService, log)
	router := server.SetupRoutes(h, log, "*")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/pokemon?limit=20&offset=0", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

