package routes

import (
	"net/http"
	"portal/internal/server/controllers"

	"github.com/gorilla/mux"
)

func Routes() http.Handler {

	r := mux.NewRouter()

	// Apply CORS middleware
	r.Use(corsMiddleware)

	baseController := controllers.NewBaseController()
	serviceController := controllers.NewServiceController(baseController)
	userController := controllers.NewUserController(baseController)

	// General routes
	r.HandleFunc("/", serviceController.HelloWorldHandler).Methods("GET")
	r.HandleFunc("/health", serviceController.HealthHandler).Methods("GET")

	// Subroute for user routes
	userRouter := r.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", userController.GetAll).Methods("GET")
	userRouter.HandleFunc("", userController.Create).Methods("POST")
	userRouter.HandleFunc("/{id}", userController.Get).Methods("GET")
	userRouter.HandleFunc("/{id}", userController.Update).Methods("PUT")
	userRouter.HandleFunc("/{id}", userController.Delete).Methods("DELETE")

	return r
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS Headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Wildcard allows all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Credentials not allowed with wildcard origins

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
