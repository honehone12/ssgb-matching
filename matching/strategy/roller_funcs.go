package strategy

import (
	"ssgb-matching/matching/matching"
	"ssgb-matching/matching/queue"
	"ssgb-matching/matching/tickets"
)

func RollerClassOnly(p matching.MatchingParams, m map[int64]*queue.Q,
) ([]matching.MatchingResult, error) {
	r := make([]matching.MatchingResult, 0)

	for class, q := range m {
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
				Class:   class,
				Matched: t,
			})
		}
	}

	return r, nil
}

func RollerClassRole(p matching.MatchingParams, m map[int64]*queue.Q,
) ([]matching.MatchingResult, error) {
	r := make([]matching.MatchingResult, 0)
	classes := int64(len(m))
	tmp := make([]int64, classes)

	for {
		total := int64(0)
		for i := int64(0); i < classes; i++ {
			l := m[i+1].Len()
			tmp[i] = l
			total += l
		}
		if total < p.MinMatchingCapacity {
			break
		}

		target := p.MinMatchingCapacity
		if total >= p.MaxMatchingCapacity {
			target = p.MaxMatchingCapacity
		}

		ts := make([]tickets.Ticket, 0, target)
		for i := int64(0); int64(len(ts)) < target; i = (i + 1) % classes {
			if tmp[i] > 0 {
				t, err := m[i+1].Dequeue()
				if err != nil {
					return nil, err
				}

				tmp[i]--
				ts = append(ts, t)
			}
		}

		r = append(r, matching.MatchingResult{
			Class:   matching.MatchingClassAny,
			Matched: ts,
		})
	}

	return r, nil
}
