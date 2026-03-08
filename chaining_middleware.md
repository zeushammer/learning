Slice pattern for Middleware Pipeline or Chain

// Instead of nesting functions, you define a custom type for your middleware and a function to "wrap" them all at once.

// Middleware is a function that wraps an http.Handler
type Middleware func(http.Handler) http.Handler

func CreateChain(middlewares ...Middleware) Middleware {
    return func(final http.Handler) http.Handler {
        // We iterate backwards so the first middleware in the slice
        // is the first one to execute in the request flow.
        for i := len(middlewares) - 1; i >= 0; i-- {
            final = middlewares[i](final)
        }
        return final
    }
}

// Here is how that cleans up your main function. It makes the execution order explicit and much easier to read:

func main() {
    mux := http.NewServeMux()
    server := NewGameServer()

    // Routes
    mux.HandleFunc("GET /player", server.game.Player)
    mux.HandleFunc("GET /queue", server.game.Queue)

    // Define the pipeline
    chain := CreateChain(
        middleware.PanicRecovery, // 1st: Catches everything
        middleware.RequestID,     // 2nd: Injects ID
        middleware.Logging,       // 3rd: Logs with ID
        middleware.RateLimit,     // 4th: Rejects early
    )

    // Wrap the mux in the chain
    handler := chain(mux)

    log.Fatal(http.ListenAndServe(":8080", handler))
}

Why this works better

    Top-to-Bottom Readability: The order in the slice is the order the request hits the server.

    Reusable Bundles: You can create different chains for different parts of your app (e.g., an ApiChain vs. an AdminChain).

    Less Boilerplate: You don't have to keep re-assigning the handler variable.

    A Note on Libraries: If you don't want to roll your own, the most popular package for this in the Go ecosystem is Alice. It does exactly this but with a slightly more polished API.


    To make this slice pattern work, your middleware functions need to follow a consistent signature: func(http.Handler) http.Handler.

A common sticking point is the RequestID middleware because it involves the context package to pass that ID down to your actual game logic.
Implementing the RequestID Middleware

Here is how you would write that specific middleware to fit into the chain we built:
Go

func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Generate or extract the ID
        id := r.Header.Get("X-Request-ID")
        if id == "" {
            id = uuid.New().String() // Or any unique string generator
        }

        // 2. Inject it into the request context
        ctx := context.WithValue(r.Context(), "trace-id", id)
        
        // 3. Set it in the response header so the client sees it
        w.Header().Set("X-Request-ID", id)

        // 4. Pass the new request (with context) to the next handler
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

How the Flow Looks

When you use the slice/chain pattern, the request travels through these layers like an onion. The RequestID middleware is a great example of "pre-processing" (setting the ID) before the "next" handler (the logger or the mux) even knows the request exists.
Using the ID in your Game Server

Since you injected the ID into the context in the middleware, your GameServer handlers can now retrieve it easily:
Go

func (s *GameServer) Player(w http.ResponseWriter, r *http.Request) {
    traceID := r.Context().Value("trace-id").(string)
    log.Printf("[%s] Fetching player data...", traceID)
    // ... logic
}

A Quick Tip on Performance

In the CreateChain function I showed you earlier, we loop backwards:
for i := len(middlewares) - 1; i >= 0; i--.

This is crucial because you want the first item in your slice to be the outermost layer of the onion. If you looped forwards, your PanicRecovery would actually end up being the very last thing to run, which defeats the purpose of catching errors from the other middlewares!

Would you like to see how to implement a Logging middleware that captures the HTTP status code (like 200 or 404), since the standard library doesn't make that obvious?