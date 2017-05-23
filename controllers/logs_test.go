package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/antonve/logger-api/controllers"
	"github.com/antonve/logger-api/utils"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func init() {
	utils.SetupTesting()
}

func TestCreateLog(t *testing.T) {
	// Setup registration request
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
	req, err := http.NewRequest(echo.POST, "/log", logBody)
	if !assert.NoError(t, err) {
		return
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controllers.APILogCreate(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, `{"success": true}`, rec.Body.String())
	}
}
