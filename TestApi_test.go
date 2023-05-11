package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddDiary(t *testing.T) {
	router := setupRouter()

	entryData := DiaryEntry{
		Title: "Test Entry",
		Body:  "This is a test diary entry.",
	}
	entryJSON, _ := json.Marshal(entryData)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/diary", bytes.NewBuffer(entryJSON))
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

}
