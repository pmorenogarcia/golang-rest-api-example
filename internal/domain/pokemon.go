package domain

import "context"

// Pokemon represents a Pokemon entity
type Pokemon struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	Height         int            `json:"height"`
	Weight         int            `json:"weight"`
	BaseExperience int            `json:"base_experience"`
	Types          []PokemonType  `json:"types"`
	Abilities      []Ability      `json:"abilities"`
	Stats          []Stat         `json:"stats"`
	Sprites        Sprites        `json:"sprites"`
}

// PokemonType represents a Pokemon type
type PokemonType struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}

// Type represents a type (like fire, water, etc.)
type Type struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Ability represents a Pokemon ability
type Ability struct {
	IsHidden bool        `json:"is_hidden"`
	Slot     int         `json:"slot"`
	Ability  AbilityInfo `json:"ability"`
}

// AbilityInfo represents ability information
type AbilityInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Stat represents a Pokemon stat
type Stat struct {
	BaseStat int      `json:"base_stat"`
	Effort   int      `json:"effort"`
	Stat     StatInfo `json:"stat"`
}

// StatInfo represents stat information
type StatInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Sprites represents Pokemon sprite images
type Sprites struct {
	FrontDefault string `json:"front_default"`
	FrontShiny   string `json:"front_shiny"`
	BackDefault  string `json:"back_default"`
	BackShiny    string `json:"back_shiny"`
}

// PokemonList represents a paginated list of Pokemon
type PokemonList struct {
	Count    int              `json:"count"`
	Next     *string          `json:"next"`
	Previous *string          `json:"previous"`
	Results  []PokemonSummary `json:"results"`
}

// PokemonSummary represents a summary of a Pokemon (used in lists)
type PokemonSummary struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// PokemonService defines the interface for Pokemon business logic
type PokemonService interface {
	// GetByName retrieves a Pokemon by name or ID
	GetByName(ctx context.Context, nameOrID string) (*Pokemon, error)

	// List retrieves a paginated list of Pokemon
	List(ctx context.Context, limit, offset int) (*PokemonList, error)
}

// PokemonClient defines the interface for external Pokemon API client
type PokemonClient interface {
	// FetchPokemon fetches a Pokemon from the external API
	FetchPokemon(ctx context.Context, nameOrID string) (*Pokemon, error)

	// FetchPokemonList fetches a list of Pokemon from the external API
	FetchPokemonList(ctx context.Context, limit, offset int) (*PokemonList, error)
}
