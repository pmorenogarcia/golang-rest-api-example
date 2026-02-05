package domain

import "errors"

var (
	// ErrPokemonNotFound is returned when a Pokemon is not found
	ErrPokemonNotFound = errors.New("pokemon not found")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")

	// ErrExternalAPI is returned when the external API fails
	ErrExternalAPI = errors.New("external API error")

)
