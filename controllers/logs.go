package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/antonve/logger-api/models"
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
)

// APILogsPost registers new log
func APILogsPost(context echo.Context) error {
	log := &models.Log{}

	// Attempt to bind request to Log struct
	err := context.Bind(log)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	user := getUser(context)
	if user == nil {
		return ServeWithError(context, 500, fmt.Errorf("could not receive user"))
	}
	log.UserID = user.ID

	// Validate request
	err = log.Validate()
	if err != nil {
		return ServeWithError(context, 400, err)
	}

	// Save to database
	logCollection := models.LogCollection{}
	_, err = logCollection.Add(log)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	return Serve(context, 201)
}

// APILogsGetAll gets all logs
func APILogsGetAll(context echo.Context) error {
	logCollection := models.LogCollection{Logs: make([]models.Log, 0)}
	user := getUser(context)
	if user == nil {
		return ServeWithError(context, 500, fmt.Errorf("could not receive user"))
	}

	// Filters
	filters := map[string]interface{}{
		"user_id":  user.ID,
		"date":     context.QueryParam("date"),
		"from":     context.QueryParam("from"),
		"until":    context.QueryParam("until"),
		"language": context.QueryParam("language"),
		"page":     context.QueryParam("page"),
	}

	err := logCollection.GetAllWithFilters(filters)

	if err != nil {
		return ServeWithError(context, 500, err)
	}

	return context.JSON(http.StatusOK, logCollection)
}

// APILogsGetByID get a single log
func APILogsGetByID(context echo.Context) error {
	logCollection := models.LogCollection{}

	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	log, err := logCollection.Get(id)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	if log == nil {
		return ServeWithError(context, 404, fmt.Errorf("no log found with id %v", id))
	}

	user := context.Get("user").(*jwt.Token).Claims.(*models.JwtClaims).User
	if !log.IsOwner(user.ID) {
		return ServeWithError(context, 403, fmt.Errorf("log doesn't belong to user"))
	}

	return context.JSON(http.StatusOK, log)
}

// APILogsUpdate updates a log
func APILogsUpdate(context echo.Context) error {
	log := &models.Log{}

	// Attempt to bind request to Log struct
	err := context.Bind(log)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	// Parse out id
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	// Update
	logCollection := models.LogCollection{}
	currentLog, err := logCollection.Get(id)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	user := getUser(context)
	if user == nil {
		return ServeWithError(context, 500, fmt.Errorf("could not receive user"))
	}

	if !currentLog.IsOwner(user.ID) {
		return ServeWithError(context, 403, fmt.Errorf("log doesn't belong to user"))
	}

	log.ID = currentLog.ID
	log.UserID = currentLog.UserID

	err = logCollection.Update(log)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	return Serve(context, 200)
}

// APILogsDelete delete a log
func APILogsDelete(context echo.Context) error {
	log := &models.Log{}

	// Parse out id
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		return ServeWithError(context, 500, err)
	}
	log.ID = id

	logCollection := models.LogCollection{}
	currentLog, err := logCollection.Get(id)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	user := getUser(context)
	if user == nil {
		return ServeWithError(context, 500, fmt.Errorf("could not receive user"))
	}

	if !currentLog.IsOwner(user.ID) {
		return ServeWithError(context, 403, fmt.Errorf("log doesn't belong to user"))
	}

	err = logCollection.Delete(log)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	return Serve(context, 200)
}
