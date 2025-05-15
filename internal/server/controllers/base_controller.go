package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"portal/internal/database"
)

type BaseController struct {
	db database.Service
}

func NewBaseController() *BaseController {
	return &BaseController{db: database.New()}
}

func (bc *BaseController) shared() {
	log.Println("this is a shared method")
}

func (bc BaseController) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}
	_, _ = w.Write(jsonResp)
}

func (bc *BaseController) HealthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(bc.db.Health())
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}
	_, _ = w.Write(jsonResp)
}
