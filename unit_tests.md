Testing middleware in Go is actually quite elegant because of the net/http/httptest package. It allows you to simulate a request/response cycle without actually opening a network port.

For the RateLimiter, we want to verify two things:

    The first request passes through with a 200 OK.

    The second request (sent immediately) is blocked with a 429 Too Many Requests.

The Rate Limit Test
Go

func TestRateLimit(t *testing.T) {
    // 1. Setup: A simple "Next" handler that just returns 200 OK
    nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    // Wrap it in our middleware
    testHandler := RateLimit(nextHandler)

    // 2. Create a mock request from a specific IP
    req := httptest.NewRequest("GET", "/", nil)
    req.RemoteAddr = "1.2.3.4:1234" // Mocking the IP

    // --- First Request (Should Pass) ---
    rr1 := httptest.NewRecorder()
    testHandler.ServeHTTP(rr1, req)

    if rr1.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", rr1.Code)
    }

    // --- Second Request (Should Fail) ---
    rr2 := httptest.NewRecorder()
    testHandler.ServeHTTP(rr2, req)

    if rr2.Code != http.StatusTooManyRequests {
        t.Errorf("Expected status 429, got %d", rr2.Code)
    }
}

Why httptest.ResponseRecorder is great

It acts like a "bucket" that catches everything the handler does. You can inspect:

    rr.Code: The HTTP status.

    rr.Body.String(): The response message.

    rr.Header(): Any headers set (like the X-Request-ID).

Pro-Tip: Testing the Panic Recovery

Testing the PanicRecovery middleware follows the same logic, but you intentionally make the "Next" handler crash:
Go

func TestPanicRecovery(t *testing.T) {
    handlerToTest := PanicRecovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        panic("something exploded!")
    }))

    req := httptest.NewRequest("GET", "/", nil)
    rr := httptest.NewRecorder()

    // This should NOT crash the test suite
    handlerToTest.ServeHTTP(rr, req)

    if rr.Code != http.StatusInternalServerError {
        t.Errorf("Expected 500 status after panic, got %d", rr.Code)
    }
}

Summary of your Middleware Architecture

You've built a system that is:

    Scalable: Easy to add new logic (Auth, Compression, etc.) by just adding to the slice.

    Testable: Each piece can be verified in isolation.

    Readable: The main function clearly describes the "laws" of your server.