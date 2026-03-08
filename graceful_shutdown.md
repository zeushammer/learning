In a production game server, you don't want to just "kill" the process while a player is in the middle of a high-stakes match. Graceful Shutdown allows your server to:

    Stop accepting new connections.

    Finish processing any active requests.

    Shut down safely once the "work" is done.

The Implementation

We use a context.WithTimeout to give the server a "deadline" (e.g., 10 seconds) to wrap up before we force it to quit.
Go

func main() {
    // 1. Setup your Mux and Middleware Chain as before
    mux := http.NewServeMux()
    handler := CreateChain(middleware.PanicRecovery, middleware.Logging)(mux)

    // 2. Define the Server
    srv := &http.Server{
        Addr:    ":8080",
        Handler: handler,
    }

    // 3. Create a channel to listen for OS signals (Ctrl+C, Kill)
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    // 4. Run the server in a goroutine so it doesn't block
    go func() {
        log.Println("Server starting on :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Listen error: %v", err)
        }
    }()

    // 5. Block here until we receive a signal
    <-stop
    log.Println("Shutting down server...")

    // 6. Create a deadline context for the shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // 7. Tell the server to shut down gracefully
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Println("Server exited cleanly")
}

Why this matters for your Game Server

    Database Integrity: If your MatchById handler is currently writing match results to a database, a hard crash could cause data corruption. This pattern ensures the SQL transaction finishes.

    Player Experience: By combining this with a "Drain" state (where you return a 503 for /queue but allow /match to finish), you can perform updates without kicking active players out of their games.

Final Architecture Review

You now have a remarkably solid Go backend. To recap:

    Middleware Slice: Cleans up the "pyramid of doom."

    RequestID: Provides a trace-id for debugging.

    Status Recorder: Allows logging of HTTP status codes.

    Panic Recovery: Prevents total crashes.

    Rate Limiting: Protects against abuse.

    CORS: Enables browser-based clients.

    Unit Tests: Proves the logic works.

    Graceful Shutdown: Protects in-flight data.