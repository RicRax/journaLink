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

	// udpate diary
	testUpdateDiary := []testCase{
		{"Update Valid Diary", `{"DID" : 1,"Body":"test"}`, http.StatusOK},
		{"Update Diary with DID 0", `{"DID" : 0, "Body": "test"}`, http.StatusBadRequest},
	}

	for _, tc := range testUpdateDiary {
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

	// delete diary

	testDeleteDiary := []testCase{
		{"Delete Valid Diary", `{"Title" : "test"}`, http.StatusOK},
		{"Delete unexisting diary", `{"Title" : ""}`, http.StatusBadRequest},
	}

	for _, tc := range testDeleteDiary {
		t.Run(tc.Name, func(t *testing.T) {
			req := httptest.NewRequest(
				"DELETE",
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
}
