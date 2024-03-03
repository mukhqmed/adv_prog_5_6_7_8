package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/verify", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function with the created request and response recorder
	VerifyHandler(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check the response body
	expectedBody := "Email parameter is missing\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}

	// Check if the redirect header is empty
	if location := rr.Header().Get("Location"); location != "" {
		t.Errorf("handler returned unexpected redirect location: got %v want empty",
			location)
	}
}
