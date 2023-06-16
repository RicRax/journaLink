package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RicRax/journalink/auth"
	"github.com/RicRax/journalink/model"
)

type testCase struct {
	Name         string
	RequestBody  string
	ExpectedCode int
}

type User = model.User

func TestAddGetUpdateDiary(t *testing.T) {
	router := setupRouter()
	auth.InitRand()
	// test adding user
	testRegister := []testCase{
		{"Register Valid user", `{"username" : "test", "password" : "test"}`, http.StatusOK},
		{
			"Register Missing Password",
			`{"username" : "test" , "password": ""}`,
			http.StatusBadRequest,
		},
		{
			"Register Missing username",
			`{"username" : "", "password":"test"}`,
			http.StatusBadRequest,
		},
	}

	for _, tc := range testRegister {
		t.Run(tc.Name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/user", bytes.NewBufferString(tc.RequestBody))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			resp := w.Result()

			if resp.StatusCode != tc.ExpectedCode {
				t.Errorf("Expected status code %d, got %d", tc.ExpectedCode, resp.StatusCode)
			}
		})
	}

	// test loggin in
	var jwt *http.Cookie

	testLogin := []testCase{
		{"Login Valid user", `{"username" : "test", "password" : "test"}`, http.StatusOK},
		{
			"Login Invaliduser",
			`{"username" : "invalid", "password" : "invalid"}`,
			http.StatusUnauthorized,
		},
		{
			"Login Missing Password",
			`{"username" : "test" , "password": ""}`,
			http.StatusUnauthorized,
		},
		{"Login Missing username", `{"username" : "", "password":"test"}`, http.StatusUnauthorized},
	}

	for _, tc := range testLogin {
		t.Run(tc.Name, func(t *testing.T) {
			req := httptest.NewRequest(
				"POST",
				"/login/authentication",
				bytes.NewBufferString(tc.RequestBody),
			)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			resp := w.Result()

			if resp.StatusCode != tc.ExpectedCode {
				t.Errorf("Expected status code %d, got %d", tc.ExpectedCode, resp.StatusCode)
			}
			if resp.StatusCode == http.StatusOK {
				jwt = resp.Cookies()[0]
			}
		})
	}

	// print(jwt)//works
	// test creating diary
	testAddDiary := []testCase{
		{"Add Valid Diary", `{"Title" : "test"}`, http.StatusOK},
		{"Add Diary with no Title", `{"Title" : ""}`, http.StatusBadRequest},
	}

	for _, tc := range testAddDiary {
		t.Run(tc.Name, func(t *testing.T) {
			req := httptest.NewRequest(
				"POST",
				"/diary",
				bytes.NewBufferString(tc.RequestBody),
			)

			req.AddCookie(jwt)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			resp := w.Result()

			if resp.StatusCode != tc.ExpectedCode {
				t.Errorf("Expected status code %d, got %d", tc.ExpectedCode, resp.StatusCode)
			}
		})
	}

	// udpate entry
	// 	modifiedEntry := DiaryInfo{
	// 		DiaryID: 1,
	// 		Title:   "Test Entry",
	// 		Body:    "Modified body",
	// 		Shared:  []string{"Riccardo", "Paolo"},
	// 	}
	// 	modifiedEntry.DiaryID = int(addedEntry.DID)
	// 	modifiedEntryJSON, _ := json.Marshal(modifiedEntry)
	// 	req, _ = http.NewRequest("POST", "/diary", bytes.NewBuffer(modifiedEntryJSON))
	// 	resp = httptest.NewRecorder()
	// 	router.ServeHTTP(resp, req)
	// 	assert.Equal(t, http.StatusOK, resp.Code)
	// }
	//
	// func TestGetSharedDiaries(t *testing.T) {
	// 	router := setupRouter()
	// 	entryData := Diary{
	// 		Title:   "Test Entry",
	// 		OwnerID: 1,
	// 		Body:    "This is a test diary entry.",
	// 	}
	// 	entryJSON, _ := json.Marshal(entryData)
	// 	resp := httptest.NewRecorder()
	// 	req, _ := http.NewRequest("POST", "/diary", bytes.NewBuffer(entryJSON))
	// 	router.ServeHTTP(resp, req)
	//
	// 	assert.Equal(t, http.StatusOK, resp.Code)
	//
	// 	// getting entry
	// 	var addedEntry Diary
	// 	json.Unmarshal(resp.Body.Bytes(), &addedEntry)
	// 	entryID := strconv.FormatUint(uint64(addedEntry.DID), 10)
	// 	req, _ = http.NewRequest("GET", "/diary/"+entryID, nil)
	// 	resp = httptest.NewRecorder()
	// 	router.ServeHTTP(resp, req)
	// 	body, _ := io.ReadAll(resp.Body)
	// 	fmt.Println(string(body))
	// 	assert.Equal(t, http.StatusOK, resp.Code)
	//
	// 	// udpate entry
	// 	modifiedEntry := DiaryInfo{
	// 		DiaryID: 1,
	// 		Title:   "Test Entry",
	// 		Body:    "Modified body",
	// 		Shared:  []string{"Riccardo", "Paolo"},
	// 	}
	// 	modifiedEntry.DiaryID = int(addedEntry.DID)
	// 	modifiedEntryJSON, _ := json.Marshal(modifiedEntry)
	// 	req, _ = http.NewRequest("POST", "/diary", bytes.NewBuffer(modifiedEntryJSON))
	// 	resp = httptest.NewRecorder()
	// 	router.ServeHTTP(resp, req)
	// 	assert.Equal(t, http.StatusOK, resp.Code)
}
