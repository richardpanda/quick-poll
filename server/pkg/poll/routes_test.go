package poll

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/richardpanda/quick-poll/server/pkg/choice"
	"github.com/richardpanda/quick-poll/server/pkg/httperror"
	"github.com/richardpanda/quick-poll/server/pkg/test"
	"github.com/stretchr/testify/assert"
)

func createChoicesTable(db *gorm.DB) {
	db.CreateTable(&choice.Choice{})
}

func createPollsTable(db *gorm.DB) {
	db.CreateTable(&Poll{})
}

func dropChoicesTable(db *gorm.DB) {
	db.DropTableIfExists("choices")
}

func dropPollsTable(db *gorm.DB) {
	db.DropTableIfExists("polls")
}

func TestPOSTPolls(t *testing.T) {
	db, close := test.DBConnection(t)
	defer close()
	createPollsTable(db)
	createChoicesTable(db)
	defer dropPollsTable(db)
	defer dropChoicesTable(db)

	b, err := json.Marshal(POSTPollsRequestBody{
		Question: "Favorite color?",
		Choices:  []string{"blue", "red", "yellow"},
	})
	assert.NoError(t, err)

	router := test.NewRouter(db, AddRoutes)
	req, err := http.NewRequest("POST", "/v1/polls", bytes.NewBuffer(b))
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var responseBody POSTPollsResponseBody
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.True(t, responseBody.ID != "")
	assert.Equal(t, "Favorite color?", responseBody.Question)
	assert.Equal(t, 3, len(responseBody.Choices))
	assert.True(t, responseBody.Choices[0].ID != "")
	assert.Equal(t, "blue", responseBody.Choices[0].Text)
}

func TestPOSTPollsWithoutRequestBody(t *testing.T) {
	db, close := test.DBConnection(t)
	defer close()

	router := test.NewRouter(db, AddRoutes)
	req, err := http.NewRequest("POST", "/v1/polls", nil)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var responseBody httperror.ResponseBody
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "Request body is missing.", responseBody.Message)
}

func TestPOSTPollsWithoutQuestion(t *testing.T) {
	db, close := test.DBConnection(t)
	defer close()

	b, err := json.Marshal(POSTPollsRequestBody{
		Choices: []string{"blue", "red", "yellow"},
	})
	assert.NoError(t, err)

	router := test.NewRouter(db, AddRoutes)
	req, err := http.NewRequest("POST", "/v1/polls", bytes.NewBuffer(b))
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var responseBody httperror.ResponseBody
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "Question is required.", responseBody.Message)
}

func TestPOSTPollsWithoutChoices(t *testing.T) {
	db, close := test.DBConnection(t)
	defer close()

	b, err := json.Marshal(POSTPollsRequestBody{
		Question: "Favorite color?",
	})
	assert.NoError(t, err)

	router := test.NewRouter(db, AddRoutes)
	req, err := http.NewRequest("POST", "/v1/polls", bytes.NewBuffer(b))
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var responseBody httperror.ResponseBody
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "Please provide at least two choices.", responseBody.Message)
}

func TestPOSTPollsWithOneChoice(t *testing.T) {
	db, close := test.DBConnection(t)
	defer close()

	b, err := json.Marshal(POSTPollsRequestBody{
		Question: "Favorite color?",
		Choices:  []string{"blue"},
	})
	assert.NoError(t, err)

	router := test.NewRouter(db, AddRoutes)
	req, err := http.NewRequest("POST", "/v1/polls", bytes.NewBuffer(b))
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var responseBody httperror.ResponseBody
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "Please provide at least two choices.", responseBody.Message)
}
