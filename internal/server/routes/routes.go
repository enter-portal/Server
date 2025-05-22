package routes

import (
	"net/http"
	"portal/internal/server/controllers"
	"portal/internal/server/middlewares"

	"github.com/gorilla/mux"
)

func Routes() http.Handler {

	r := mux.NewRouter()

	// Create CORS middleware
	corsMiddleware := middlewares.NewCorsMiddleware()
	canonicalPathMiddleware := middlewares.NewCanonicalPathMiddleware()
	// Define a map of paths and their corresponding HTTP methods which are allowed without authentication
	authMiddleware := middlewares.NewAuthMiddleware(map[string][]string{
		"/":       {"GET"},
		"/health": {"GET"},
		"/users":  {"GET", "POST"},
	})

	// Middlewares
	r.Use(corsMiddleware.Middleware)
	r.Use(authMiddleware.Middleware)

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

	return canonicalPathMiddleware.Middleware(r)
}
