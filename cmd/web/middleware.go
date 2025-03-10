package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	
	"github.com/justinas/nosurf"
)

const (
	authenticatedUserIDSessionManager = "authenticated_user_id"
	userRoleSessionManager            = "user_role"
)

// commonHeaders middleware sets common HTTP headers and generates a nonce for script security.
//
// Parameters:
//
//	next - The next handler in the chain
//
// Returns:
//
//	http.Handler - A new handler that includes the common headers and nonce
func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		// generating a nonce for the script embedded in the templates
		nonce, err := newNonce()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		
		// putting the nonce in the context
		ctx := context.WithValue(r.Context(), nonceContextKey, nonce)
		r = r.WithContext(ctx)
		
		// setting the common headers
		w.Header().Set("Content-Security-Policy", fmt.Sprintf("script-src 'self' 'nonce-%s' https://cdn.jsdelivr.net https://unpkg.com/lucide@latest", nonce))
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		
		w.Header().Set("Server", "Golang server")
		
		next.ServeHTTP(w, r)
	})
}

// logRequest middleware logs incoming HTTP requests.
//
// Parameters:
//
//	next - The next handler in the chain
//
// Returns:
//
//	http.Handler - A new handler that logs the request details
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)
		
		// DEBUG
		app.logger.Debug("received request", slog.String("ip", ip), slog.String("protocol", proto), slog.String("method", method), slog.String("URI", uri))
		
		next.ServeHTTP(w, r)
	})
}

// recoverPanic middleware recovers from any panics that occur during the handling of a request.
//
// Parameters:
//
//	next - The next handler in the chain
//
// Returns:
//
//	http.Handler - A new handler that recovers from panics and logs errors
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("recovering from panic: %s", err))
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}

// noSurf middleware adds CSRF protection to the application.
//
// Parameters:
//
//	next - The next handler in the chain
//
// Returns:
//
//	http.Handler - A new handler with CSRF protection enabled
func noSurf(next http.Handler) http.Handler {
	
	csrfHandler := nosurf.New(next)
	
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	
	return csrfHandler
}
