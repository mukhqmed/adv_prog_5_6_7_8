package auth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	// Create a request to the LoginHandler
	form := url.Values{}
	form.Set("username", "validUsername")
	form.Set("password", "validPassword")
	req, err := http.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	LoginHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}

	// Check the Location header for the redirection URL
	expectedLocation := "/admin" // Change this to the expected redirection URL
	if location := rr.Header().Get("Location"); location != expectedLocation {
		t.Errorf("handler returned unexpected redirect location: got %v want %v",
			location, expectedLocation)
	}
}
