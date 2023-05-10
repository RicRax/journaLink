package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDiaryEntries(t *testing.T) {
	// Create a new test database
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&DiaryEntry{})
	defer db.Migrator().DropTable(&DiaryEntry{})

	// Create a new Gin router and set up the routes
	r := gin.Default()
	r.POST("/diary", func(c *gin.Context) {
		AddDiaryEntry(c, db)
	})
	r.GET("/diary/:id", func(c *gin.Context) {
		GetDiaryEntry(c, db)
	})
	r.PUT("/diary/:id", func(c *gin.Context) {
		UpdateDiaryEntry(c, db)
	})
	r.GET("/diary", func(c *gin.Context) {
		GetAllDiaryEntries(c, db)
	})

	// Test adding a new diary entry
	entryData := DiaryEntry{
		Title: "Test Entry",
		Body:  "This is a test diary entry.",
	}
	entryJSON, _ := json.Marshal(entryData)
	req, _ := http.NewRequest("POST", "/diary", bytes.NewBuffer(entryJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Test getting the added diary entry
	var addedEntry DiaryEntry
	json.Unmarshal(resp.Body.Bytes(), &addedEntry)
	entryID := addedEntry.ID
	req, _ = http.NewRequest("GET", "/diary/"+entryID, nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Test updating the diary entry
	updatedEntryData := DiaryEntry{
		Title: "Updated Test Entry",
		Body:  "This is an updated test diary entry.",
	}
	updatedEntryJSON, _ := json.Marshal(updatedEntryData)
	req, _ = http.NewRequest("PUT", "/diary/"+entryID, bytes.NewBuffer(updatedEntryJSON))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Test getting all diary entries
	req, _ = http.NewRequest("GET", "/diary", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	var entries []DiaryEntry
	json.Unmarshal(resp.Body.Bytes(), &entries)
	assert.Equal(t, 1, len(entries))
}
