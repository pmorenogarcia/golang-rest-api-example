package domain

// TypeEffectiveness maps Pokemon types to types they are strong against
var TypeEffectiveness = map[string][]string{
	"normal":   {},
	"fire":     {"grass", "ice", "bug", "steel"},
	"water":    {"fire", "ground", "rock"},
	"electric": {"water", "flying"},
	"grass":    {"water", "ground", "rock"},
	"ice":      {"grass", "ground", "flying", "dragon"},
	"fighting": {"normal", "ice", "rock", "dark", "steel"},
	"poison":   {"grass", "fairy"},
	"ground":   {"fire", "electric", "poison", "rock", "steel"},
	"flying":   {"grass", "fighting", "bug"},
	"psychic":  {"fighting", "poison"},
	"bug":      {"grass", "psychic", "dark"},
	"rock":     {"fire", "ice", "flying", "bug"},
	"ghost":    {"psychic", "ghost"},
	"dragon":   {"dragon"},
	"dark":     {"psychic", "ghost"},
	"steel":    {"ice", "rock", "fairy"},
	"fairy":    {"fighting", "dragon", "dark"},
}

// IsStrongAgainst checks if attackType is strong against defenseType
func IsStrongAgainst(attackType, defenseType string) bool {
	strongAgainst, exists := TypeEffectiveness[attackType]
	if !exists {
		return false
	}

	for _, t := range strongAgainst {
		if t == defenseType {
			return true
		}
	}
	return false
}