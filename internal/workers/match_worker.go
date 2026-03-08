package workers

import (
	"log"

	"learning/internal/models"

	"github.com/google/uuid"
)

func StartMatchmaker(
	queue <-chan models.MatchRequest,
	matches chan<- models.Match,
) {

	go func() {

		log.Println("matchmaker started")

		var buffer []models.MatchRequest

		for req := range queue {

			buffer = append(buffer, req)

			log.Println("buffer size:", len(buffer))

			if len(buffer) >= 2 {

				match := models.Match{
					ID: uuid.NewString(),
					Players: []models.MatchRequest{
						buffer[0],
						buffer[1],
					},
				}

				matches <- match

				buffer = buffer[2:]
			}
		}
	}()
}
