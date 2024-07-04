package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func addSecurityHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		h.ServeHTTP(w, r)
	})
}

func main() {
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 10, // in seconds
		HttpOnly: true,
		Secure:   true, // set HTTPS-only cookies
		SameSite: http.SameSiteStrictMode,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", authMiddleware(indexHandler))
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/new", authMiddleware(uploadHandler))
	mux.HandleFunc("/download/", authMiddleware(fileHandler))

	http.Handle("/", addSecurityHeaders(mux))

	fmt.Println("Starting server at https://localhost:8443")
	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
