package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestAddGetUpdateDiary(t *testing.T) {
	router := setupRouter()

	//adding diary
	entryData := DiaryEntry{
		Title: "Test Entry",
		Body:  "This is a test diary entry.",
	}
	entryJSON, _ := json.Marshal(entryData)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/diary", bytes.NewBuffer(entryJSON))
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	//getting entry
	// var addedEntry DiaryEntry
	entryID := strconv.FormatUint(uint64(entryData.ID), 10)
	req, _ = http.NewRequest("GET", "/diary/"+entryID, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	assert.Equal(t, http.StatusOK, resp.Code)

	//
}
