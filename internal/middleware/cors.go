package middleware

import (
	"net/http"
)

var allowedOrigins = map[string]bool{
	"https://blog-steel-pi-43.vercel.app":      true,
	"https://blog-ajif8sjjr-acmev2.vercel.app": true,
	"https://blog-three-indol-95.vercel.app":   true,
}

// CORSMiddleware adiciona cabeçalhos CORS para permitir requisições de origens específicas.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
