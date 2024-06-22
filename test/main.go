package main

import (
	"encoding/json"
	"net/http"

	"github.com/tvative/goapitest"
)

var instance *apitest.Instance

func main() {
	// Initialize instance
	instance = apitest.Initialize(&apitest.Config{
		Level:        apitest.VerboseLevel,
		IsNeedResult: true,
	}, true)

	// Assign handlers
	instance.Mux.HandleFunc("/hello", HelloHandler)

	// Test cases
	testCases := []apitest.TestCase{
		{
			ID:             "TC_01",
			Type:           apitest.HappyPath,
			Details:        "Sample case",
			EndPoint:       "/hello",
			Method:         http.MethodGet,
			Header:         nil,
			BodyType:       "",
			QueryParams:    nil,
			BodyParams:     nil,
			ExpectedResult: "",
			ExpectedStatus: http.StatusOK,
			IsIgnored:      false,
		}, {
			ID:             "TC_01",
			Type:           apitest.HappyPath,
			Details:        "Another sample case",
			EndPoint:       "/hello",
			Method:         http.MethodGet,
			Header:         nil,
			BodyType:       "",
			QueryParams:    nil,
			BodyParams:     nil,
			ExpectedResult: "",
			ExpectedStatus: http.StatusBadRequest,
			IsIgnored:      false,
		},
	}

	// Add test cases
	for _, tc := range testCases {
		err := instance.Add(tc)
		if err != nil {
			panic(err)
		}
	}

	// Dump test cases
	instance.Cases.Dump(instance)
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "hello world!"}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		panic(err)
	}
}
