package matchmaking

import (
	"sort"

	"learning/internal/models"

	"github.com/google/uuid"
)

func MatchBySkill(reqs []models.MatchRequest) []models.Match {

	sort.Slice(reqs, func(i, j int) bool {
		return reqs[i].Rating < reqs[j].Rating
	})

	var matches []models.Match

	for len(reqs) >= 2 {

		m := models.Match{
			ID: uuid.NewString(),
			Players: []models.MatchRequest{
				reqs[0],
				reqs[1],
			},
		}

		matches = append(matches, m)

		reqs = reqs[2:]
	}

	return matches
}
