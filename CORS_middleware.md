The CORS Middleware

CORS works in two steps:

    The Preflight: The browser sends an OPTIONS request to see if it's allowed to talk to you.

    The Actual Request: If the preflight passes, it sends your GET or POST.

Go

func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set the allowed origin (use "*" for development, specific URL for production)
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Request-ID")

        // Handle the preflight request
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}

The Final "Pro-Level" Middleware Stack

Now that we have all the pieces, let's look at the final assembly in your main.go. This setup is robust, safe, and incredibly easy to maintain.
Go

func main() {
    mux := http.NewServeMux()
    server := NewGameServer()

    // Routes
    mux.HandleFunc("GET /player", server.game.Player)
    mux.HandleFunc("GET /match/{id}", server.game.MatchById)

    // THE PIPELINE
    // Order: Outside -> Inside
    // 1. Recovery: Catch crashes
    // 2. CORS: Handle browser security
    // 3. RateLimit: Block spammers early
    // 4. RequestID: ID the "survivors"
    // 5. Logging: Record what happened
    finalHandler := CreateChain(
        middleware.PanicRecovery,
        middleware.CORS,
        middleware.RateLimit,
        middleware.RequestID,
        middleware.Logging,
    )(mux)

    log.Printf("Game Server starting on :%s", SERVERPORT)
    log.Fatal(http.ListenAndServe(":"+SERVERPORT, finalHandler))
}

Why this is "Clean"

    Decoupled: Your GameServer logic doesn't know anything about CORS, Logging, or Panics. It just focuses on game logic.

    Fail-Fast: The RateLimit is near the top. If a user is spamming, they get rejected before the server spends any CPU cycles generating a RequestID or writing a log entry.

    Traceability: Because Logging is the innermost middleware, it has access to the RequestID set just above it, and it can record the Status Code set by any middleware or handler inside it.