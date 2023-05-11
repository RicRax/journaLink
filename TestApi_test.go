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
	var addedEntry DiaryEntry
	json.Unmarshal(resp.Body.Bytes(), &addedEntry)
	entryID := strconv.FormatUint(uint64(addedEntry.ID), 10)
	req, _ = http.NewRequest("GET", "/diary/"+entryID, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	assert.Equal(t, http.StatusOK, resp.Code)

	//udpate entry
	modifiedEntry := DiaryEntry{Title: "Modified", Body: "Modified body"}
	modifiedEntryJSON, _ := json.Marshal(modifiedEntry)
	req, _ = http.NewRequest("PUT", "/diary/"+entryID, bytes.NewBuffer(modifiedEntryJSON))
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	//
}
