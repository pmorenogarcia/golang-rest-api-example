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

// PokemonCount represents the count of Pokemon
type PokemonCount struct {
	Count int `json:"count"`
}

// PokemonComparison represents the result of comparing two Pokemon
type PokemonComparison struct {
	Message       string   `json:"message"`
	TypeAdvantage string   `json:"type_advantage,omitempty"`
	Winner        *Pokemon `json:"winner"`
	Pokemon1      *Pokemon `json:"pokemon1"`
	Pokemon2      *Pokemon `json:"pokemon2"`
}

// PokemonService defines the interface for Pokemon business logic
type PokemonService interface {
	// GetByName retrieves a Pokemon by name or ID
	GetByName(ctx context.Context, nameOrID string) (*Pokemon, error)

	// GetCount retrieves the total count of Pokemon
	GetCount(ctx context.Context) (*PokemonCount, error)

	// ComparePokemon compares two Pokemon based on type effectiveness
	ComparePokemon(ctx context.Context, name1, name2 string) (*PokemonComparison, error)
}

// PokemonClient defines the interface for external Pokemon API client
type PokemonClient interface {
	// FetchPokemon fetches a Pokemon from the external API
	FetchPokemon(ctx context.Context, nameOrID string) (*Pokemon, error)

	// FetchPokemonCount fetches the total count from the external API
	FetchPokemonCount(ctx context.Context) (int, error)
}
