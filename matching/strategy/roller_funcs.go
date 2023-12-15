package strategy

import (
	"ssgb-matching/errs"
	"ssgb-matching/matching/matching"
	"ssgb-matching/matching/queue"
)

func RollerClassOnly(p matching.MatchingParams, m map[int64]*queue.Q,
) ([]matching.MatchingResult, error) {
	r := make([]matching.MatchingResult, 0)

	for _, q := range m {
		for l := q.Len(); l >= p.MinMatchingCapacity; l = q.Len() {
			n := p.MinMatchingCapacity
			if l >= p.MaxMatchingCapacity {
				n = p.MaxMatchingCapacity
			}

			t, err := q.DequeueN(n)
			if err != nil {
				return nil, err
			}

			r = append(r, matching.MatchingResult{
				Matched: t,
			})
		}
	}

	return r, nil
}

func RollerClassRole(p matching.MatchingParams, m map[int64]*queue.Q,
) ([]matching.MatchingResult, error) {
	return nil, errs.ErrorNotimplemented()
}
