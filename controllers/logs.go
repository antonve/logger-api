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
		return Return500(context, err)
	}

	user := getUser(context)
	if user == nil {
		return Return500(context, fmt.Errorf("could not receive user"))
	}
	log.UserID = user.ID

	// Validate request
	err = log.Validate()
	if err != nil {
		return Return400(context, err)
	}

	// Save to database
	logCollection := models.LogCollection{}
	_, err = logCollection.Add(log)
	if err != nil {
		return Return500(context, err)
	}

	return Return201(context)
}

// APILogsGetAll gets all logs
func APILogsGetAll(context echo.Context) error {
	logCollection := models.LogCollection{Logs: make([]models.Log, 0)}
	user := getUser(context)
	if user == nil {
		return Return500(context, fmt.Errorf("could not receive user"))
	}

	err := logCollection.GetAllFromUser(user.ID)

	if err != nil {
		return Return500(context, err)
	}

	return context.JSON(http.StatusOK, logCollection)
}

// APILogsGetByID get a single log
func APILogsGetByID(context echo.Context) error {
	logCollection := models.LogCollection{}

	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		return Return500(context, err)
	}

	log, err := logCollection.Get(id)
	if err != nil {
		return Return500(context, err)
	}

	if log == nil {
		return Return404(context, fmt.Errorf("no log found with id %v", id))
	}

	user := context.Get("user").(*jwt.Token).Claims.(*models.JwtClaims).User
	if !log.IsOwner(user.ID) {
		return Return403(context, fmt.Errorf("log doesn't belong to user"))
	}

	return context.JSON(http.StatusOK, log)
}

// APILogsUpdate updates a log
func APILogsUpdate(context echo.Context) error {
	log := &models.Log{}

	// Attempt to bind request to Log struct
	err := context.Bind(log)
	if err != nil {
		return Return500(context, err)
	}

	// Parse out id
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		return Return500(context, err)
	}

	// Update
	logCollection := models.LogCollection{}
	currentLog, err := logCollection.Get(id)
	if err != nil {
		return Return500(context, err)
	}

	user := getUser(context)
	if user == nil {
		return Return500(context, fmt.Errorf("could not receive user"))
	}

	if !currentLog.IsOwner(user.ID) {
		return Return403(context, fmt.Errorf("log doesn't belong to user"))
	}

	err = logCollection.Update(log)
	if err != nil {
		return Return500(context, err)
	}

	return Return200(context)
}

// APILogsDelete delete a log
func APILogsDelete(context echo.Context) error {
	log := &models.Log{}

	// Parse out id
	fmt.Println(context.Param("id"))
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		return Return500(context, err)
	}
	log.ID = id

	logCollection := models.LogCollection{}
	currentLog, err := logCollection.Get(id)
	if err != nil {
		return Return500(context, err)
	}

	user := getUser(context)
	if user == nil {
		return Return500(context, fmt.Errorf("could not receive user"))
	}

	if !currentLog.IsOwner(user.ID) {
		return Return403(context, fmt.Errorf("log doesn't belong to user"))
	}

	err = logCollection.Delete(log)
	if err != nil {
		return Return500(context, err)
	}

	return Return200(context)
}
