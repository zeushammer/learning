package matchmaking

import (
	"learning/internal/models"

	"github.com/google/uuid"
)

func MatchPlayers(reqs []models.MatchRequest) []models.Match {

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
