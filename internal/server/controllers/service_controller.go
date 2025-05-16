package controllers

import (
	"net/http"
	"portal/internal/database"
	"portal/internal/server/utils"
)

type ServiceController struct {
	*BaseController
}

func NewServiceController(baseController *BaseController) *ServiceController {
	return &ServiceController{baseController}
}

func (sc *ServiceController) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSONResponse(w, http.StatusOK, utils.JSON{"message": "Hello World"})
}

func (sc *ServiceController) HealthHandler(w http.ResponseWriter, r *http.Request) {
	httpStatusCode, json := database.Health()
	utils.WriteJSONResponse(w, httpStatusCode, json)
}
