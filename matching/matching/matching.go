package matching

import (
	"ssgb-matching/matching/queue"
	"ssgb-matching/matching/tickets"
)

type RollerFunc func(MatchingParams, map[int64]*queue.Q,
) ([]MatchingResult, error)

type MatchingParams struct {
	MinMatchingCapacity int64
	MaxMatchingCapacity int64
}

type MatchingResult struct {
	Matched []tickets.Ticket
}
