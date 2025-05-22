package middlewares

import (
	"net/http"
	"portal/internal/server/utils"
	"strings"
)

// AuthMiddleware middleware structure
type AuthMiddleware struct {
	avoidAuthPaths map[string][]string
}

// NewAuthMiddleware creates a new instance of the AuthMiddleware
func NewAuthMiddleware(avoidAuthPaths map[string][]string) *AuthMiddleware {
	return &AuthMiddleware{
		avoidAuthPaths: avoidAuthPaths,
	}
}

// Helper function to check if a string exists in a slice
func (am *AuthMiddleware) contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// AuthMiddleware validates the JWT token and user ID before processing the request.
// It checks the Authorization header for the token and x-key header for the user ID.
func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// IMPORTANT: Normalize the path before checking avoidAuthPaths
		// This assumes a path normalization middleware (like CanonicalPathMiddleware)
		// has already run and modified r.URL.Path if necessary.
		requestPath := r.URL.Path
		// Check if the path is allowed for the request method without authentication
		allowedMethods, pathExists := am.avoidAuthPaths[requestPath]
		// TODO: Add support for wildcards
		if pathExists {
			if am.contains(allowedMethods, r.Method) {
				next.ServeHTTP(w, r)
				return
			}
		}
		// Extract token from Authorization header or x-access-token header
		authorizationHeader := r.Header.Get("Authorization")
		var tokenString string
		if authorizationHeader != "" {
			parts := strings.Split(authorizationHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" { // Ensure it's a Bearer token
				tokenString = parts[1]
			}
		}
		// If not found in Authorization header, try x-access-token
		if tokenString == "" {
			tokenString = r.Header.Get("x-access-token")
		}
		// Extract user ID from x-key header
		userId := r.Header.Get("x-key")
		// Validate the token and user ID
		if tokenString != "" && userId != "" {
			result, err := utils.ValidateToken(tokenString, userId)
			if err != nil {
				// Log the error for debugging purposes (optional but recommended)
				// log.Printf("Token validation error: %v", err)
				utils.WriteJSONResponse(w, http.StatusUnauthorized, utils.JSON{"message": "Authentication failed: " + err.Error()})
				return
			}
			if result {
				// Token is valid, proceed to the next handler
				next.ServeHTTP(w, r)
			} else {
				// Token is invalid but no specific error from ValidateToken
				utils.WriteJSONResponse(w, http.StatusUnauthorized, utils.JSON{"message": "Invalid Token"})
			}
		} else {
			// Missing token or key
			utils.WriteJSONResponse(w, http.StatusUnauthorized, utils.JSON{"message": "Authentication required: Token or Key is missing"})
		}
	})
}
