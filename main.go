// main.go
package main

import (
	"context"
	"net/http"
	"os"

	playGo "learning/game"
	"learning/internal/middleware"
	"learning/internal/utils"

	"learning/internal/logger"
	"learning/internal/requestctx"
)

// docker run <container> command args
// go run main.go command args

var SERVERPORT = "8080"

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
	server := NewGameServer()

	mux.HandleFunc("GET /player", server.game.Player)
	mux.HandleFunc("GET /queue", server.game.Queue)
	mux.HandleFunc("GET /match/{id}", server.game.MatchById)

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

	// Wrap ListenAndServer
	err := http.ListenAndServe("localhost:"+SERVERPORT, handler)

	if err != nil {
		logger.Error(ctx, err, "Server failed to start or shut down unexpectedly")
		// todo: capture SIGNAL, graceful shutdown, also graceful if error
		os.Exit(1)
	}

}

type gameServer struct {
	game      *playGo.Game
	multitool utils.Leatherman
}

func NewGameServer() *gameServer {
	ticTacGame := playGo.New()
	return &gameServer{game: ticTacGame}
}
