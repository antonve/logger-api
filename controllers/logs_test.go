package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/antonve/logger-api/controllers"
	"github.com/antonve/logger-api/models"
	"github.com/antonve/logger-api/models/enums"
	"github.com/antonve/logger-api/utils"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func init() {
	utils.SetupTesting()
}

func TestCreateLog(t *testing.T) {
	// Setup create log request
	e := echo.New()
	logBody := strings.NewReader(`{
    "language": "JA",
    "date": "2017-05-23",
    "duration": 25,
    "activity": "READING",
    "notes": {
      "type": "BOOK",
      "series": "キングダム",
      "volume": 1,
      "pages": 200
    }
  }`)
	req, err := http.NewRequest(echo.POST, "/api/logs", logBody)
	if !assert.NoError(t, err) {
		return
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controllers.APILogsPost(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, `{"success": true}`, rec.Body.String())
	}
}

func TestGetLog(t *testing.T) {
	// Setup log to grab
	log := models.Log{Language: enums.LanguageKorean, Date: "2016-04-05", Duration: 60, Activity: enums.ActivityListening}
	logCollection := models.LogCollection{}
	id, _ := logCollection.Add(&log)

	// Setup log request
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/logs/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))

	if assert.NoError(t, controllers.APILogsGetByID(c)) {
		// Check response
		var body models.Log
		assert.Equal(t, http.StatusOK, rec.Code)
		err := json.Unmarshal(rec.Body.Bytes(), &body)

		// Check if the log has information
		assert.Nil(t, err)
		assert.NotEmpty(t, body.Language)
		assert.NotEmpty(t, body.Date)
		assert.NotEmpty(t, body.Duration)
		assert.NotEmpty(t, body.Activity)

		// Check if the log has the correct information
		assert.Equal(t, enums.LanguageKorean, body.Language)
		assert.Equal(t, "2016-04-05", body.Date)
		assert.Equal(t, uint64(60), body.Duration)
		assert.Equal(t, enums.ActivityListening, body.Activity)
	}
}
