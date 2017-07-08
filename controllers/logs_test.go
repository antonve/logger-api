package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/antonve/logger-api/config"
	"github.com/antonve/logger-api/controllers"
	"github.com/antonve/logger-api/models"
	"github.com/antonve/logger-api/models/enums"
	"github.com/antonve/logger-api/utils"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
)

type LogsBody struct {
	Logs []models.Log `json:"logs"`
}

var mockJwtToken string
var mockUser *models.User

func init() {
	utils.SetupTesting()
	mockJwtToken, mockUser = utils.SetupTestUser()
}

func TestPost(t *testing.T) {
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mockJwtToken))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, middleware.JWTWithConfig(config.GetJWTConfig(&models.JwtClaims{}))(controllers.APILogsPost)(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, `{"success": true}`, rec.Body.String())
	}
}

func TestGetByID(t *testing.T) {
	// Setup log to grab
	log := models.Log{UserID: mockUser.ID, Language: enums.LanguageKorean, Date: "2016-10-05", Duration: 60, Activity: enums.ActivityListening}
	logCollection := models.LogCollection{}
	id, _ := logCollection.Add(&log)

	// Setup log request
	e := echo.New()
	req := httptest.NewRequest(echo.GET, fmt.Sprintf("/api/logs/%d", id), nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mockJwtToken))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/logs/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))

	if assert.NoError(t, middleware.JWTWithConfig(config.GetJWTConfig(&models.JwtClaims{}))(controllers.APILogsGetByID)(c)) {
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
		assert.Equal(t, "2016-10-05", body.Date)
		assert.Equal(t, uint64(60), body.Duration)
		assert.Equal(t, enums.ActivityListening, body.Activity)
	}
}

func TestGetAll(t *testing.T) {
	// Setup log to grab
	logCollection := models.LogCollection{}
	var ids [3]uint64
	ids[0], _ = logCollection.Add(&models.Log{UserID: mockUser.ID, Language: enums.LanguageJapanese, Date: "2016-04-04", Duration: 30, Activity: enums.ActivityGrammar})
	ids[1], _ = logCollection.Add(&models.Log{UserID: mockUser.ID, Language: enums.LanguageMandarin, Date: "2016-04-03", Duration: 45, Activity: enums.ActivityOther})
	ids[2], _ = logCollection.Add(&models.Log{UserID: mockUser.ID, Language: enums.LanguageKorean, Date: "2016-04-05", Duration: 55, Activity: enums.ActivityTextbook})

	// Setup log request
	e := echo.New()

	req := httptest.NewRequest(echo.GET, "/api/logs", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mockJwtToken))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/logs")

	// Without any filters
	if assert.NoError(t, middleware.JWTWithConfig(config.GetJWTConfig(&models.JwtClaims{}))(controllers.APILogsGetAll)(c)) {
		// Check response
		var body LogsBody
		assert.Equal(t, http.StatusOK, rec.Code)
		err := json.Unmarshal(rec.Body.Bytes(), &body)

		// Check if the log has information
		assert.Nil(t, err)
		assert.True(t, len(body.Logs) >= 3)
	}

	// By date
	req = httptest.NewRequest(echo.GET, "/api/logs?date=2016-04-03", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mockJwtToken))

	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/api/logs")

	if assert.NoError(t, middleware.JWTWithConfig(config.GetJWTConfig(&models.JwtClaims{}))(controllers.APILogsGetAll)(c)) {
		// Check response
		var body LogsBody
		assert.Equal(t, http.StatusOK, rec.Code)
		err := json.Unmarshal(rec.Body.Bytes(), &body)

		// Check if the log has information
		assert.Nil(t, err)
		assert.True(t, len(body.Logs) == 1)
	}

	// By daterange
	req = httptest.NewRequest(echo.GET, "/api/logs?from=2016-04-01&until=2016-04-04", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mockJwtToken))

	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/api/logs")

	if assert.NoError(t, middleware.JWTWithConfig(config.GetJWTConfig(&models.JwtClaims{}))(controllers.APILogsGetAll)(c)) {
		// Check response
		var body LogsBody
		assert.Equal(t, http.StatusOK, rec.Code)
		err := json.Unmarshal(rec.Body.Bytes(), &body)

		// Check if the log has information
		assert.Nil(t, err)
		assert.True(t, len(body.Logs) == 2)
	}
}

func TestUpdate(t *testing.T) {
	// Setup log to grab
	logCollection := models.LogCollection{}
	id, _ := logCollection.Add(&models.Log{UserID: mockUser.ID, Language: enums.LanguageGerman, Date: "2016-03-30", Duration: 5, Activity: enums.ActivityTranslation})

	// Setup log request
	e := echo.New()
	logBody := strings.NewReader(`{
    "language": "KR",
    "date": "2017-03-30",
    "duration": 25,
    "activity": "GRAMMAR"
  }`)
	req := httptest.NewRequest(echo.PUT, fmt.Sprintf("/api/logs/%d", id), logBody)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mockJwtToken))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/logs/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))

	if assert.NoError(t, middleware.JWTWithConfig(config.GetJWTConfig(&models.JwtClaims{}))(controllers.APILogsUpdate)(c)) {
		// Check response
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"success": true}`, rec.Body.String())

		log, _ := logCollection.Get(id)
		assert.Equal(t, enums.LanguageKorean, log.Language)
		assert.Equal(t, "2017-03-30", log.Date)
		assert.Equal(t, uint64(25), log.Duration)
		assert.Equal(t, enums.ActivityGrammar, log.Activity)
	}
}

func TestDelete(t *testing.T) {
	// Setup log to grab
	logCollection := models.LogCollection{}
	id, _ := logCollection.Add(&models.Log{UserID: mockUser.ID, Language: enums.LanguageJapanese, Date: "2016-01-30", Duration: 50, Activity: enums.ActivityFlashcards})

	// Setup log request
	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, fmt.Sprintf("/api/logs/%d", id), nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mockJwtToken))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/logs/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))

	if assert.NoError(t, middleware.JWTWithConfig(config.GetJWTConfig(&models.JwtClaims{}))(controllers.APILogsDelete)(c)) {
		// Check response
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"success": true}`, rec.Body.String())

		log, e := logCollection.Get(id)
		fmt.Println(e)
		assert.Nil(t, log)
	}
}
