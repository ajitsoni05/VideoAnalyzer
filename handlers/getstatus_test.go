package handlers

import (
	"fmt"
	"github/poornachandra7707/myboilerplate/models"

	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var testInputs = []models.GetStatusTestInputs{
	{Key: "hulu_70", ExpectedStatus: "Status of video:failed"},
	{Key: "ytb_323", ExpectedStatus: "Status of video:completed"},
	{Key: "Sling_23", ExpectedStatus: "Status of video:failed"},
	{Key: "hulu_65", ExpectedStatus: "Status of video:completed"},
}

func TestGetStatus(t *testing.T) {

	// Start a local HTTP server
	router := mux.NewRouter()
	router.HandleFunc("/getstatus/{key}", GetStatus).Methods("GET")
	server := httptest.NewServer(router)
	defer server.Close()

	for ind, _ := range testInputs {
		// Prepare a request to send to the server
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/getstatus/%s", server.URL, testInputs[ind].Key), nil)

		// Send the request to the server
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("Error sending the request: %v", err)
		}
		// Assert the response status
		if res.StatusCode != http.StatusOK {
			t.Errorf("Unexpected status code: got %v want %v", res.StatusCode, http.StatusOK)
		}

		// Assert the response body
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		body := string(bodyBytes)
		if body != testInputs[ind].ExpectedStatus {
			t.Errorf("Unexpected response body: got %v want %v", body, testInputs[ind].ExpectedStatus)
		}
	}
}
