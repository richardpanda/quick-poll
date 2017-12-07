package router_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/richardpanda/quick-poll/server/pkg/choice"
	"github.com/richardpanda/quick-poll/server/pkg/httperror"
	"github.com/richardpanda/quick-poll/server/pkg/poll"
	. "github.com/richardpanda/quick-poll/server/pkg/router"
	"github.com/richardpanda/quick-poll/server/pkg/test"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetPoll(t *testing.T) {
	db, close := test.DBConnection(t)
	defer close()
	test.CreatePollsTable(db)
	test.CreateChoicesTable(db)
	defer test.DropPollsTable(db)
	defer test.DropChoicesTable(db)

	p := poll.Poll{
		ID:       uuid.NewV4().String(),
		Question: "Favorite color?",
		Choices: []choice.Choice{
			choice.New("blue"),
			choice.New("red"),
			choice.New("yellow"),
		},
	}

	err := db.Create(&p).Error
	assert.NoError(t, err)

	router := NewTestRouter(db)
	endpoint := fmt.Sprintf("/v1/polls/%s", p.ID)
	req, err := http.NewRequest("GET", endpoint, nil)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var responseBody poll.GETPollResponseBody
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.True(t, responseBody.ID != "")
	assert.Equal(t, "Favorite color?", responseBody.Question)
	assert.Equal(t, 3, len(responseBody.Choices))
	assert.True(t, responseBody.Choices[0].ID != "")
	assert.Equal(t, "blue", responseBody.Choices[0].Text)
	assert.Equal(t, 0, responseBody.Choices[0].NumVotes)
}

func TestGetPollWithInvalidID(t *testing.T) {
	db, close := test.DBConnection(t)
	defer close()
	test.CreatePollsTable(db)
	test.CreateChoicesTable(db)
	defer test.DropPollsTable(db)
	defer test.DropChoicesTable(db)

	p := poll.Poll{
		ID:       uuid.NewV4().String(),
		Question: "Favorite color?",
		Choices: []choice.Choice{
			choice.New("blue"),
			choice.New("red"),
			choice.New("yellow"),
		},
	}

	err := db.Create(&p).Error
	assert.NoError(t, err)

	invalidID := uuid.NewV4()
	router := NewTestRouter(db)
	endpoint := fmt.Sprintf("/v1/polls/%s", invalidID)
	req, err := http.NewRequest("GET", endpoint, nil)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var responseBody httperror.ResponseBody
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "Invalid poll ID.", responseBody.Message)
}

func TestPOSTPolls(t *testing.T) {
	db, close := test.DBConnection(t)
	defer close()
	test.CreatePollsTable(db)
	test.CreateChoicesTable(db)
	defer test.DropPollsTable(db)
	defer test.DropChoicesTable(db)

	b, err := json.Marshal(poll.POSTPollsRequestBody{
		Question: "Favorite color?",
		Choices:  []string{"blue", "red", "yellow"},
	})
	assert.NoError(t, err)

	router := NewTestRouter(db)
	req, err := http.NewRequest("POST", "/v1/polls", bytes.NewBuffer(b))
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var responseBody poll.POSTPollsResponseBody
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.True(t, responseBody.ID != "")
	assert.Equal(t, "Favorite color?", responseBody.Question)
	assert.Equal(t, 3, len(responseBody.Choices))
	assert.True(t, responseBody.Choices[0].ID != "")
	assert.Equal(t, "blue", responseBody.Choices[0].Text)
	assert.Equal(t, 0, responseBody.Choices[0].NumVotes)
}

func TestPOSTPollsWithoutRequestBody(t *testing.T) {
	db, close := test.DBConnection(t)
	defer close()

	router := NewTestRouter(db)
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

	b, err := json.Marshal(poll.POSTPollsRequestBody{
		Choices: []string{"blue", "red", "yellow"},
	})
	assert.NoError(t, err)

	router := NewTestRouter(db)
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

	b, err := json.Marshal(poll.POSTPollsRequestBody{
		Question: "Favorite color?",
	})
	assert.NoError(t, err)

	router := NewTestRouter(db)
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

	b, err := json.Marshal(poll.POSTPollsRequestBody{
		Question: "Favorite color?",
		Choices:  []string{"blue"},
	})
	assert.NoError(t, err)

	router := NewTestRouter(db)
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
