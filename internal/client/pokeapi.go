package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/polgarcia/golang-rest-api/internal/domain"
	"github.com/polgarcia/golang-rest-api/pkg/logger"
	"go.uber.org/zap"
)

const (
	maxRetries     = 3
	retryDelay     = 1 * time.Second
	retryBackoff   = 2.0
)

// PokeAPIClient implements the PokemonClient interface
type PokeAPIClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *logger.Logger
}

// NewPokeAPIClient creates a new PokeAPI client
func NewPokeAPIClient(baseURL string, timeout time.Duration, log *logger.Logger) *PokeAPIClient {
	return &PokeAPIClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
		logger: log,
	}
}

// FetchPokemon fetches a Pokemon from the PokeAPI
func (c *PokeAPIClient) FetchPokemon(ctx context.Context, nameOrID string) (*domain.Pokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%s", c.baseURL, strings.ToLower(nameOrID))

	c.logger.Debug("Fetching Pokemon",
		zap.String("name_or_id", nameOrID),
		zap.String("url", url),
	)

	var pokemon domain.Pokemon
	if err := c.doRequestWithRetry(ctx, url, &pokemon); err != nil {
		if err == domain.ErrPokemonNotFound {
			c.logger.Debug("Pokemon not found", zap.String("name_or_id", nameOrID))
			return nil, domain.ErrPokemonNotFound
		}
		c.logger.Error("Failed to fetch Pokemon",
			zap.String("name_or_id", nameOrID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("%w: %v", domain.ErrExternalAPI, err)
	}

	c.logger.Info("Successfully fetched Pokemon",
		zap.String("name", pokemon.Name),
		zap.Int("id", pokemon.ID),
	)

	return &pokemon, nil
}

// doRequestWithRetry performs an HTTP request with retry logic
func (c *PokeAPIClient) doRequestWithRetry(ctx context.Context, url string, result interface{}) error {
	var lastErr error
	delay := retryDelay

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Debug("Retrying request",
				zap.Int("attempt", attempt+1),
				zap.Int("max_retries", maxRetries),
				zap.Duration("delay", delay),
			)

			// Wait before retrying
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}

			// Exponential backoff
			delay = time.Duration(float64(delay) * retryBackoff)
		}

		err := c.doRequest(ctx, url, result)
		if err == nil {
			return nil
		}

		// Don't retry on 404
		if err == domain.ErrPokemonNotFound {
			return err
		}

		lastErr = err

		// Don't retry on context cancellation
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	c.logger.Warn("Max retries exceeded", zap.Error(lastErr))
	return lastErr
}

// doRequest performs a single HTTP GET request
func (c *PokeAPIClient) doRequest(ctx context.Context, url string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "golang-rest-api/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Handle HTTP errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Debug("HTTP error",
			zap.Int("status_code", resp.StatusCode),
			zap.String("body", string(body)),
		)

		if resp.StatusCode == http.StatusNotFound {
			return domain.ErrPokemonNotFound
		}

		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Decode response
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
