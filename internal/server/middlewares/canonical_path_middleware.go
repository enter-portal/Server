package middlewares

import (
	"log"
	"net/http"
	"strings"
)

type CanonicalPathMiddleware struct{}

func NewCanonicalPathMiddleware() *CanonicalPathMiddleware {
	return &CanonicalPathMiddleware{}
}

// TrimSuffixSlash middleware
func (cpm *CanonicalPathMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// remove the trailing slash from URL Path if it's not root
		log.Println(r.RequestURI)
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}
