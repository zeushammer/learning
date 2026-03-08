// main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"learning/internal/handlers"
	"learning/internal/middleware"
	"learning/internal/models"
	"learning/internal/workers"

	"learning/internal/logger"
	"learning/internal/requestctx"

	sdk "agones.dev/agones/sdks/go"
)

// docker run <container> command args
// go run main.go command args

var SERVERPORT = "1337"

var matchQueue = make(chan models.MatchRequest, 100)

func main() {

	/* request logging
	   simple rate limiting
	   panic recovery
	   context with trace-id
	*/

	// Create a root context for the application lifetime
	ctx := context.Background()
	ctx = requestctx.SetTraceID(ctx, "SYSTEM-STARTUP")
	logger.Info(ctx, "Initializing game server...")

	// Setup the mux
	mux := http.NewServeMux()

	handlers.Queue = matchQueue

	//mux.HandleFunc("GET /player", server.game.Player)
	mux.HandleFunc("POST /queue", handlers.QueueHandler)
	//mux.HandleFunc("GET /match/{id}", server.game.MatchById)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("game server ready"))
	})

	// Define the pipeline
	// Top-to-Bottom Readability: The order in the slice is the order the request hits the server.
	chain := middleware.CreateChain(
		middleware.PanicRecovery, // 1st: Catches everything
		middleware.RequestID,     // 2nd: Injects ID
		middleware.Logging,       // 3rd: Logs with ID
		// middleware.RateLimit,     // 4th: Rejects early
		// todo: RateLimit
	)

	// Wrap the mux in the chain
	handler := chain(mux)

	// Remember this is essentially doing: handler becomes PanicRecovery(RequestID(Logging(RateLimit(mux)))).
	// request --> [Panic Recovery] --> [RequestID] --> [Logging] --> [RateLimit] --> [Mux] --> [Handler]

	logger.Info(ctx, "Server starting on %s", SERVERPORT)

	// queue where players enter matchmaking
	matchQueue := make(chan models.MatchRequest, 100)

	// channel where completed matches appear
	matchResults := make(chan models.Match, 100)

	// give the HTTP handler access to the queue
	handlers.Queue = matchQueue

	// start worker pool
	workers.StartMatchmaker(
		matchQueue,   // input queue
		matchResults, // output matches
	)

	// listen for match results
	go func() {

		for match := range matchResults {

			log.Printf(
				"match created %s players=%v",
				match.ID,
				match.Players,
			)

			// later this is where you would allocate
			// a game server with Agones
		}

	}()

	// connect to Agones
	s, err := sdk.NewSDK()
	if err != nil {
		logger.Error(ctx, err, "Could not connect to Agones SDK")
	}

	// tell Agones this server is ready to accept players
	if err := s.Ready(); err != nil {
		logger.Error(ctx, err, "Cloud not mark Ready")
	}

	// start health pings
	go func() {
		for {
			if err := s.Health(); err != nil {
				logger.Error(ctx, err, "health ping failed")
			}
			time.Sleep(2 * time.Second)
		}
	}()

	// Wrap ListenAndServer
	err = http.ListenAndServe(":"+SERVERPORT, handler)

	if err != nil {
		logger.Error(ctx, err, "Server failed to start or shut down unexpectedly")
		// todo: capture SIGNAL, graceful shutdown, also graceful if error
		os.Exit(1)
	}

}
