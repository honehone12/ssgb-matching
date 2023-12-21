package strategy

import (
	"errors"
	"ssgb-matching/matching/matching"
)

const (
	StrategyClassOnly int = 1 + iota
	StrategyClassRole
)

func StrategyFunc(strategy int) (matching.RollerFunc, error) {
	switch strategy {
	case StrategyClassOnly:
		return RollerClassOnly, nil
	case StrategyClassRole:
		return RollerClassRole, nil
	default:
		return nil, errors.New("no roller func is defined to the strategy")
	}
}
