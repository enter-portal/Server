package controllers

import (
	"io"
	"net/http"
)

type UserController struct {
	*BaseController
}

func NewUserController(baseController *BaseController) *UserController {
	return &UserController{baseController}
}

// Create user
func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Create\n")
}

// Update user
func (uc *UserController) Update(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Update\n")
}

// Delete user
func (uc *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Delete\n")
}

// Get user
func (uc *UserController) Get(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Get\n")
}

// Get all users
func (uc *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "GetAll\n")
}
