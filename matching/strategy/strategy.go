package strategy

import (
	"errors"
	"ssgb-matching/matching/matching"
)

const (
	StrategyClassOnly int = 1 + iota
	StrategyClassRole
)

func ValidStrategy(strategy int) error {
	switch strategy {
	case StrategyClassOnly:
		fallthrough
	case StrategyClassRole:
		return nil
	default:
		return errors.New("invalid matching strategy")
	}
}

func StrategyFunc(strategy int) (matching.RollerFunc, error) {
	switch strategy {
	case StrategyClassOnly:
		return RollerClassOnly, nil
	case StrategyClassRole:
		return RollerClassRole, nil
	default:
		return nil, errors.New("no RollerFunc is defined to the strategy")
	}
}
