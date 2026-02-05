package domain

import "errors"

var (
	// ErrPokemonNotFound is returned when a Pokemon is not found
	ErrPokemonNotFound = errors.New("pokemon not found")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")

	// ErrExternalAPI is returned when the external API fails
	ErrExternalAPI = errors.New("external API error")

	// ErrInvalidLimit is returned when the limit parameter is invalid
	ErrInvalidLimit = errors.New("invalid limit: must be between 1 and 100")

	// ErrInvalidOffset is returned when the offset parameter is invalid
	ErrInvalidOffset = errors.New("invalid offset: must be non-negative")
)
