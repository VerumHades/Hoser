package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/internal/database"
	"server/internal/handlers"
	"testing"
)

var app = handlers.NewTestApp()

func testHandlerForUser(user database.User, request *http.Request, handler http.HandlerFunc) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()

	session, _ := app.Store.Get(request, "user-session")
	session.Values["authenticated"] = true
	session.Values["user"] = user
	_ = session.Save(request, response)

	handler.ServeHTTP(response, request)

	return response
}

func testRequestAndExpectRentals(user database.User, t *testing.T) []handlers.ApiUserRental {
	request := httptest.NewRequest("GET", "/user/rentals", nil)
	response := testHandlerForUser(user, request, app.UserRentalsRequestHandler())

	if response.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", response.Code)
	}

	var rentals []handlers.ApiUserRental
	if err := json.NewDecoder(response.Body).Decode(&rentals); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	return rentals
}

func testRentingRequest(user database.User, t *testing.T, body map[string]string) {
	bodyBytes, _ := json.Marshal(body)

	request := httptest.NewRequest("POST", "/user/rental", bytes.NewReader(bodyBytes))
	request.Header.Set("Content-Type", "application/json")

	response := testHandlerForUser(user, request, app.UserRentRequestHandler())

	_, has_id := body["id"]
	if !has_id && response.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", response.Code)
	} else {
		return
	}

	if response.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", response.Code)
	}
}
func TestUserRentalsRequestHandler(t *testing.T) {
	dummyUserAlice, _ := app.DatabaseInteractor.GetUserByName("alice")
	dummyUserBob, _ := app.DatabaseInteractor.GetUserByName("bob")

	rentals := testRequestAndExpectRentals(dummyUserAlice, t)
	if len(rentals) != 1 {
		t.Errorf("expected at least one rental for alice")
	}
	if rentals[0].Title != "Go API Testbed" {
		t.Errorf("expected at rental name to be 'Go API Testbed'")
	}

	rentals = testRequestAndExpectRentals(dummyUserBob, t)
	if len(rentals) != 1 {
		t.Errorf("expected at least one rental for bob")
	}
	if rentals[0].Title != "TensorFlow Training Node" {
		t.Errorf("expected at rental name to be 'TensorFlow Training Node'")
	}
}

func TestUserRentRequestHandler(t *testing.T) {
	dummyUserAlice, _ := app.DatabaseInteractor.GetUserByName("alice")
	dummyUserBob, _ := app.DatabaseInteractor.GetUserByName("bob")

	testRentingRequest(dummyUserAlice, t, map[string]string{})
	testRentingRequest(dummyUserAlice, t, map[string]string{"id": "0"})

	testRentingRequest(dummyUserBob, t, map[string]string{"id": "1"})
	testRentingRequest(dummyUserBob, t, map[string]string{"id": "1"})

	aliceRentals, _ := dummyUserAlice.GetRentals()
	bobRentals, _ := dummyUserBob.GetRentals()

	if len(aliceRentals) != 2 {
		t.Errorf("expected at two rentals for alice")
	}
	if aliceRentals[1].GetTitle() != "Node.js Dev Stack" {
		t.Errorf("expected alices new rental title to be: 'New Rental'")
	}

	if len(bobRentals) != 3 {
		t.Errorf("expected at two rentals for bob")
	}
	if bobRentals[1].GetDescription() != "Go + Kafka setup" {
		t.Errorf("expected bobs new rental description to be: 'New Rental Description'")
	}

	if bobRentals[2].GetTitle() != "Go Microservices Boilerplate" {
		t.Errorf("expected bobs new unset rental title to be: 'New Rental'")
	}
	if bobRentals[2].GetDescription() != "Go + Kafka setup" {
		t.Errorf("expected bobs new unset rental description to be: 'New Rental Description'")
	}
}
