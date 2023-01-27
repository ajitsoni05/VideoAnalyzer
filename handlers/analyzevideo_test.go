package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	//"github/poornachandra7707/myboilerplate/handlers"
	"github/poornachandra7707/myboilerplate/models"

	"github.com/gorilla/mux"
)

var AnalyzeVideotestInputs = []models.AnalyzeVideoTestInputs{
	{AnalyzeVideoReqBody: models.AnalyzeVideoReqBody{PlatformName: "hulu", AssetID: "70", FilePath: "htps://brightcove-test-delivery.s3.amazonaws.com/itv_source.mov"}, ExpectedStatus: "Please note down the asset key for further reference:ucplay_123\n"},
	{AnalyzeVideoReqBody: models.AnalyzeVideoReqBody{PlatformName: "ytb", AssetID: "323", FilePath: "https://brightcove-test-delivery.s3.amazonaws.com/itv_source.mov"}, ExpectedStatus: "Status of video:in progress\n"},
	{AnalyzeVideoReqBody: models.AnalyzeVideoReqBody{PlatformName: "Sling", AssetID: "23", FilePath: "htps://brightcove-test-delivery.s3.amazonaws.com/itv_source.mov"}, ExpectedStatus: "Please note down the asset key for further reference:ucvoot_21\n"},
	{AnalyzeVideoReqBody: models.AnalyzeVideoReqBody{PlatformName: "hulu", AssetID: "65", FilePath: "htts://brightcove-test-delivery.s3.amazonaws.com/itv_source.mov"}, ExpectedStatus: "Status of video:Already done previously\n"},
}

func TestAnalyzeVideo(t *testing.T) {
	// Start a local HTTP server
	router := mux.NewRouter()
	router.HandleFunc("/analyze", AnalyzeVideo).Methods("POST")
	server := httptest.NewServer(router)
	defer server.Close()
	// Prepare a request body to send to the server
	// var analysisResult = models.AnalysisResult{
	//  PlatformName: "test_platform",
	//  AssetID:      "test_asset_id",
	//  FilePath:     "https://brightcove-test-delivery.s3.amazonaws.com/itv_source.mov",
	// }
	for ind, _ := range AnalyzeVideotestInputs {
		body, _ := json.Marshal(AnalyzeVideotestInputs[ind].AnalyzeVideoReqBody)
		// Prepare a request to send to the server
		req, _ := http.NewRequest("POST", fmt.Sprintf("%s/analyze", server.URL), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
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
		// var returnedAnalysisResult models.AnalysisResult
		// json.NewDecoder(res.Body).Decode(&returnedAnalysisResult)
		// if analysisResult != returnedAnalysisResult {
		//  t.Errorf("Unexpected response body: got %v want %v", returnedAnalysisResult, analysisResult)
		// }
		// bodyBytes, _ := ioutil.ReadAll(res.Body)
		// body = string(bodyBytes)
		// if body = AnalyzeVideotestInputs[ind].ExpectedStatus {
		//  t.Errorf("Unexpected response body: got %v want %v", body, AnalyzeVideotestInputs[ind].ExpectedStatus)
		// }
		resBody, _ := ioutil.ReadAll(res.Body)
		returnedAnalysisResult := string(resBody)
		if returnedAnalysisResult != AnalyzeVideotestInputs[ind].ExpectedStatus {
			t.Errorf("Unexpected response body: got %v want %v", returnedAnalysisResult, AnalyzeVideotestInputs[ind].ExpectedStatus)
		}
	}
}
