package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/antonve/logger-api/models"

	"github.com/labstack/echo"
)

// APILogRegister registers new log
func APILogCreate(context echo.Context) error {
	log := &models.Log{}

	// Attempt to bind request to Log struct
	err := context.Bind(log)
	if err != nil {
		return Return500(context, err)
	}

	// Validate request
	err = log.Validate()
	if err != nil {
		return Return400(context, err)
	}

	// Save to database
	logCollection := models.LogCollection{}
	err = logCollection.Add(log)
	if err != nil {
		return Return500(context, err)
	}

	return Return201(context)
}

// APILogGetAll gets all logs
func APILogGetAll(context echo.Context) error {
	logCollection := models.LogCollection{Logs: make([]models.Log, 0)}
	err := logCollection.GetAll()

	if err != nil {
		return Return500(context, err)
	}

	return context.JSON(http.StatusOK, logCollection)
}

// APILogGetByID get a single log
func APILogGetByID(context echo.Context) error {
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
		return Return404(context, fmt.Errorf("No Log found with id %v", id))
	}

	return context.JSON(http.StatusOK, log)
}

// APILogUpdate updates a log
func APILogUpdate(context echo.Context) error {
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
	log.ID = id

	// Update
	logCollection := models.LogCollection{}
	err = logCollection.Update(log)
	if err != nil {
		return Return500(context, err)
	}

	return Return201(context)
}
