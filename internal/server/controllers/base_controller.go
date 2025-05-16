package controllers

import (
	"portal/internal/database"

	"gorm.io/gorm"
)

type BaseController struct {
	db *gorm.DB
}

func NewBaseController() *BaseController {
	return &BaseController{db: database.New()}
}
