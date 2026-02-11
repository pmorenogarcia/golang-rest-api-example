package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/polgarcia/golang-rest-api/internal/domain"
	"github.com/polgarcia/golang-rest-api/pkg/logger"
	"go.uber.org/zap"
)

// PokemonService implements the domain.PokemonService interface
type PokemonService struct {
	client domain.PokemonClient
	logger *logger.Logger
}

// NewPokemonService creates a new Pokemon service
func NewPokemonService(client domain.PokemonClient, log *logger.Logger) *PokemonService {
	return &PokemonService{
		client: client,
		logger: log,
	}
}

// GetByName retrieves a Pokemon by name or ID
func (s *PokemonService) GetByName(ctx context.Context, nameOrID string) (*domain.Pokemon, error) {
	// Validate input
	if nameOrID == "" {
		s.logger.Debug("Invalid input: empty name or ID")
		return nil, fmt.Errorf("%w: name or ID cannot be empty", domain.ErrInvalidInput)
	}

	// Normalize name to lowercase
	nameOrID = strings.ToLower(strings.TrimSpace(nameOrID))

	s.logger.Info("Getting Pokemon",
		zap.String("name_or_id", nameOrID),
	)

	// Fetch Pokemon from client
	pokemon, err := s.client.FetchPokemon(ctx, nameOrID)
	if err != nil {
		s.logger.Error("Failed to get Pokemon",
			zap.String("name_or_id", nameOrID),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Info("Successfully retrieved Pokemon",
		zap.String("name", pokemon.Name),
		zap.Int("id", pokemon.ID),
	)

	return pokemon, nil
}

// GetCount retrieves the total count of Pokemon
func (s *PokemonService) GetCount(ctx context.Context) (*domain.PokemonCount, error) {
	// Fetch count from client
	count, err := s.client.FetchPokemonCount(ctx)
	if err != nil {
		s.logger.Error("Failed to get Pokemon count", zap.Error(err))
		return nil, err
	}

	return &domain.PokemonCount{Count: count}, nil
}

// ComparePokemon compares two Pokemon based on type effectiveness
func (s *PokemonService) ComparePokemon(ctx context.Context, name1, name2 string) (*domain.PokemonComparison, error) {
	// Validate inputs
	if name1 == "" || name2 == "" {
		s.logger.Debug("Invalid input: empty Pokemon names")
		return nil, fmt.Errorf("%w: both Pokemon names must be provided", domain.ErrInvalidInput)
	}

	// Normalize names
	name1 = strings.ToLower(strings.TrimSpace(name1))
	name2 = strings.ToLower(strings.TrimSpace(name2))

	if name1 == name2 {
		s.logger.Debug("Invalid input: same Pokemon names")
		return nil, fmt.Errorf("%w: cannot compare a Pokemon with itself", domain.ErrInvalidInput)
	}

	s.logger.Info("Comparing Pokemon",
		zap.String("pokemon1", name1),
		zap.String("pokemon2", name2),
	)

	// Fetch both Pokemon
	pokemon1, err := s.client.FetchPokemon(ctx, name1)
	if err != nil {
		s.logger.Error("Failed to fetch first Pokemon",
			zap.String("name", name1),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to fetch %s: %w", name1, err)
	}

	pokemon2, err := s.client.FetchPokemon(ctx, name2)
	if err != nil {
		s.logger.Error("Failed to fetch second Pokemon",
			zap.String("name", name2),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to fetch %s: %w", name2, err)
	}

	// Compare types and determine winner
	comparison := s.compareTypes(pokemon1, pokemon2)

	winnerName := "tie"
	if comparison.Winner != nil {
		winnerName = comparison.Winner.Name
	}

	s.logger.Info("Pokemon comparison completed",
		zap.String("winner", winnerName),
		zap.String("message", comparison.Message),
	)

	return comparison, nil
}

// compareTypes compares the types of two Pokemon and determines the winner
func (s *PokemonService) compareTypes(p1, p2 *domain.Pokemon) *domain.PokemonComparison {
	// Get primary types for neutral matchup display
	var primaryType1, primaryType2 string
	if len(p1.Types) > 0 {
		primaryType1 = p1.Types[0].Type.Name
	}
	if len(p2.Types) > 0 {
		primaryType2 = p2.Types[0].Type.Name
	}

	// Track type advantages
	var p1AdvantageType, p1AdvantageAgainst string
	var p2AdvantageType, p2AdvantageAgainst string
	p1HasAdvantage := false
	p2HasAdvantage := false

	// Check all type combinations and track which types create advantages
	for _, t1 := range p1.Types {
		for _, t2 := range p2.Types {
			if domain.IsStrongAgainst(t1.Type.Name, t2.Type.Name) {
				p1HasAdvantage = true
				p1AdvantageType = t1.Type.Name
				p1AdvantageAgainst = t2.Type.Name
			}
			if domain.IsStrongAgainst(t2.Type.Name, t1.Type.Name) {
				p2HasAdvantage = true
				p2AdvantageType = t2.Type.Name
				p2AdvantageAgainst = t1.Type.Name
			}
		}
	}

	comparison := &domain.PokemonComparison{
		Pokemon1: p1,
		Pokemon2: p2,
	}

	// Determine winner and message
	if p1HasAdvantage && !p2HasAdvantage {
		comparison.Winner = p1
		comparison.Message = fmt.Sprintf("%s is stronger than %s because %s type beats %s type",
			capitalize(p1.Name), capitalize(p2.Name), p1AdvantageType, p1AdvantageAgainst)
		comparison.TypeAdvantage = fmt.Sprintf("%s beats %s", p1AdvantageType, p1AdvantageAgainst)
	} else if p2HasAdvantage && !p1HasAdvantage {
		comparison.Winner = p2
		comparison.Message = fmt.Sprintf("%s is stronger than %s because %s type beats %s type",
			capitalize(p2.Name), capitalize(p1.Name), p2AdvantageType, p2AdvantageAgainst)
		comparison.TypeAdvantage = fmt.Sprintf("%s beats %s", p2AdvantageType, p2AdvantageAgainst)
	} else if p1HasAdvantage && p2HasAdvantage {
		// Both have advantages, it's a tie
		comparison.Winner = nil
		comparison.Message = fmt.Sprintf("%s and %s both have type advantages against each other - it's a tie!",
			capitalize(p1.Name), capitalize(p2.Name))
		comparison.TypeAdvantage = fmt.Sprintf("%s beats %s, but %s beats %s",
			p1AdvantageType, p1AdvantageAgainst, p2AdvantageType, p2AdvantageAgainst)
	} else {
		// No type advantage, neutral matchup
		comparison.Winner = nil
		comparison.Message = fmt.Sprintf("Neither %s (%s) nor %s (%s) has a type advantage - it's a neutral matchup!",
			capitalize(p1.Name), primaryType1, capitalize(p2.Name), primaryType2)
	}

	return comparison
}

// capitalize capitalizes the first letter of a string
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
