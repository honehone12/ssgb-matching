package matching

import (
	"ssgb-matching/matching/queue"
	"ssgb-matching/matching/tickets"
)

const (
	MatchingClassAny = 0
)

type RollerFunc func(MatchingParams, map[int64]*queue.Q,
) ([]MatchingResult, error)

type MatchingParams struct {
	MinMatchingCapacity int64
	MaxMatchingCapacity int64
}

type MatchingResult struct {
	Class   int64
	Matched []tickets.Ticket
}
