package controllers

import (
	"net/http"
	"portal/internal/server/models"
	"portal/internal/server/utils"

	"github.com/gorilla/mux"
)

// UserController handles requests related to users
type UserController struct {
	*BaseController
}

// NewUserController creates a new UserController instance
func NewUserController(baseController *BaseController) *UserController {
	return &UserController{baseController}
}

// Create user
func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Parse the request body into the user struct
	if err := utils.ParseJSONRequestBody(r.Body, &user); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, utils.JSON{"message": err.Error()})
		return
	}

	// Check if user already exists with the same email
	var existingUser models.User
	if err := uc.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		utils.WriteJSONResponse(w, http.StatusConflict, utils.JSON{"message": "User with this email already exists"})
		return
	}

	// Check if user already exists with the same username
	if err := uc.db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		utils.WriteJSONResponse(w, http.StatusConflict, utils.JSON{"message": "User with this username already exists"})
		return
	}

	// If no user exists with the same email or username, create the new user
	if err := uc.db.Create(&user).Error; err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.JSON{"message": "Error creating user"})
		return
	}

	// Return success response
	utils.WriteJSONResponse(w, http.StatusCreated, utils.JSON{"message": "User created successfully", "user": user.ToJSON()})
}

// Update user
func (uc *UserController) Update(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the path
	vars := mux.Vars(r)
	id := vars["id"]

	var user models.User

	// Parse the request body into the user struct
	if err := utils.ParseJSONRequestBody(r.Body, &user); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, utils.JSON{"message": err.Error()})
		return
	}

	// Check if user exists
	var existingUser models.User
	if err := uc.db.First(&existingUser, id).Error; err != nil {
		utils.WriteJSONResponse(w, http.StatusNotFound, utils.JSON{"message": "User not found"})
		return
	}

	// Update the user
	if err := uc.db.Model(&existingUser).Updates(user).Error; err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.JSON{"message": err.Error()})
		return
	}

	// Return success response
	utils.WriteJSONResponse(w, http.StatusOK, utils.JSON{"message": "User updated successfully", "user": existingUser.ToJSON()})
}

// Delete user
func (uc *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the path
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if user exists
	var user models.User
	if err := uc.db.First(&user, id).Error; err != nil {
		utils.WriteJSONResponse(w, http.StatusNotFound, utils.JSON{"message": "User not found"})
		return
	}

	// Delete the user
	if err := uc.db.Delete(&user).Error; err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.JSON{"message": "Error deleting user"})
		return
	}

	// Return success response
	utils.WriteJSONResponse(w, http.StatusOK, utils.JSON{"message": "User deleted successfully"})
}

// Get user
func (uc *UserController) Get(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the path
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if user exists
	var user models.User
	if err := uc.db.First(&user, id).Error; err != nil {
		utils.WriteJSONResponse(w, http.StatusNotFound, utils.JSON{"message": "User not found"})
		return
	}

	// Return the user data
	utils.WriteJSONResponse(w, http.StatusOK, utils.JSON{"user": user.ToJSON()})
}

// Get all users
func (uc *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	// Query to get all users
	if err := uc.db.Find(&users).Error; err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.JSON{"message": "Error retrieving users"})
		return
	}

	// Convert users to JSON format
	userList := make([]utils.JSON, len(users))
	for i, user := range users {
		userList[i] = user.ToJSON()
	}

	// Return the user data
	utils.WriteJSONResponse(w, http.StatusOK, utils.JSON{"users": userList})
}
