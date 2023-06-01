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
	entryData := Diary{
		Title:   "Test Entry",
		OwnerID: 1, Body: "This is a test diary entry.",
	}
	entryJSON, _ := json.Marshal(entryData)
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/diary", bytes.NewBuffer(entryJSON))
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	//getting entry
	var addedEntry Diary
	json.Unmarshal(resp.Body.Bytes(), &addedEntry)
	entryID := strconv.FormatUint(uint64(addedEntry.ID), 10)
	req, _ = http.NewRequest("GET", "/diary/"+entryID, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	assert.Equal(t, http.StatusOK, resp.Code)

	//udpate entry
	modifiedEntry := DiaryInfo{DiaryID: 1, Title: "Test Entry", Body: "Modified body", Shared: []string{"Riccardo", "Paolo"}}
	modifiedEntry.DiaryID = int(addedEntry.ID)
	modifiedEntryJSON, _ := json.Marshal(modifiedEntry)
	req, _ = http.NewRequest("POST", "/diary", bytes.NewBuffer(modifiedEntryJSON))
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetSharedDiaries(t *testing.T) {
	router := setupRouter()
	entryData := Diary{
		Title:   "Test Entry",
		OwnerID: 1,
		Body:    "This is a test diary entry.",
	}
	entryJSON, _ := json.Marshal(entryData)
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/diary", bytes.NewBuffer(entryJSON))
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	//getting entry
	var addedEntry Diary
	json.Unmarshal(resp.Body.Bytes(), &addedEntry)
	entryID := strconv.FormatUint(uint64(addedEntry.ID), 10)
	req, _ = http.NewRequest("GET", "/diary/"+entryID, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	assert.Equal(t, http.StatusOK, resp.Code)

	//udpate entry
	modifiedEntry := DiaryInfo{DiaryID: 1, Title: "Test Entry", Body: "Modified body", Shared: []string{"Riccardo", "Paolo"}}
	modifiedEntry.DiaryID = int(addedEntry.ID)
	modifiedEntryJSON, _ := json.Marshal(modifiedEntry)
	req, _ = http.NewRequest("POST", "/diary", bytes.NewBuffer(modifiedEntryJSON))
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
