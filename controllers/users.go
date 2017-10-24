package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/antonve/logger-api/models"
	"github.com/antonve/logger-api/models/enums"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// APIUserGetAll gets all users
func APIUserGetAll(context echo.Context) error {
	userCollection := models.UserCollection{Users: make([]models.User, 0)}
	err := userCollection.GetAll()

	if err != nil {
		return ServeWithError(context, 500, err)
	}

	return context.JSON(http.StatusOK, userCollection)
}

// APIUserGetByID get the profile of a user
func APIUserGetByID(context echo.Context) error {
	userCollection := models.UserCollection{}

	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	currentUser := context.Get("currentUser").(*jwt.Token).Claims.(*models.JwtClaims).User
	if !(currentUser != nil && (currentUser.ID == id || currentUser.Role == enums.RoleAdmin)) {
		return ServeWithError(context, 403, fmt.Errorf("not allowed to access this user"))
	}

	user, err := userCollection.Get(id)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	if user == nil {
		return ServeWithError(context, 404, fmt.Errorf("no User found with id %v", id))
	}

	return context.JSON(http.StatusOK, user)
}

// APIUserUpdate updates a user
func APIUserUpdate(context echo.Context) error {
	user := &models.User{}

	// Attempt to bind request to User struct
	err := context.Bind(user)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	// Parse out id
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		return ServeWithError(context, 500, err)
	}
	user.ID = id

	currentUser := getUser(context)
	if !(currentUser != nil && (currentUser.ID == id || currentUser.Role == enums.RoleAdmin)) {
		return ServeWithError(context, 403, fmt.Errorf("not allowed to access this user"))
	}

	// Return an error when an unauthorized request to change role was made
	if user.Role != enums.RoleAdmin && user.Role != "" {
		return ServeWithError(context, 403, fmt.Errorf("not allowed to change the user role"))
	}

	if user.Role == "" {
		user.Role = currentUser.Role
	}

	// Update
	userCollection := models.UserCollection{}
	err = userCollection.Update(user)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	return Serve(context, 200)
}
