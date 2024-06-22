// GO API Testing Package
//
// Copyright (c) 2024 Tvative.
// All rights reserved.
//
// Source code and its usage are governed by
// the MIT license.

package apitest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

// TestCaseType is the type of test case
type TestCaseType int

const (
	HappyPath       TestCaseType = iota // HappyPath is the happy path test case
	EdgeCase                            // EdgeCase is the edge case test case
	NegativeCase                        // NegativeCase is the negative case test case
	BoundaryCase                        // BoundaryCase is the boundary case test case
	CornerCase                          // CornerCase is the corner case test case
	StressCase                          // StressCase is the stress case test case
	SmokeCase                           // SmokeCase is the smoke case test case
	RegressionCase                      // RegressionCase is the regression case test case
	IntegrationCase                     // IntegrationCase is the integration case test case
)

// TestCase is the struct for test case
type TestCase struct {
	ID             string         // ID is the test case ID
	Type           TestCaseType   // Type is the test case type
	Details        string         // Details is the test case details
	EndPoint       string         // EndPoint is the test case API endpoint
	Method         string         // Method is the test case method
	Header         map[string]any // Header is the test case header
	BodyType       string         // BodyType is the test case body type
	QueryParams    any            // QueryParams is the test case query parameters
	BodyParams     any            // BodyParams is the test case body parameters
	ExpectedResult string         // ExpectedResult is the test case expected result
	ExpectedStatus int            // ExpectedStatus is the test case expected status code
	IsIgnored      bool           // IsIgnored is the flag to ignore the test case
}

// TestCases is the struct for test cases
type TestCases struct {
	ResultGot           string               // ResultGot is the test case result
	StatusGot           string               // StatusGot is the test case status
	StatusCodeGot       int                  // StatusCodeGot is the test case status code
	ProtoGot            string               // ProtoGot is the test case protocol
	ProtoMajorGot       int                  // ProtoMajorGot is the test case protocol major
	ProtoMinorGot       int                  // ProtoMinorGot is the test case protocol minor
	ContentLengthGot    int64                // ContentLengthGot is the test case content length
	TransferEncodingGot []string             // TransferEncodingGot is the test case transfer encoding
	IsUncompressed      bool                 // IsUncompressed is the flag to check if the test case is uncompressed
	TLSGot              *tls.ConnectionState // TLSGot is the test case TLS connection state
	Time                time.Duration        // Time is the test case time
	Case                TestCase             // Case is the test case
	Next                *TestCases           // Next is the next test case
}

// insert is the function to insert test case result
func (h *TestCases) insert(result TestCases) error {
	if result.Case.IsIgnored {
		return nil
	}

	newNode := &result
	if h == nil {
		*h = *newNode
		return nil
	}

	current := h
	for current.Next != nil {
		current = current.Next
	}

	current.Next = newNode
	return nil
}

// isEmptyNode is the function to check if the test case is empty
func isEmptyNode(tc TestCases) bool {
	return tc.ResultGot == "" && tc.StatusGot == "" && tc.StatusCodeGot == 0 &&
		tc.ProtoGot == "" && tc.ProtoMajorGot == 0 && tc.ProtoMinorGot == 0 &&
		tc.ContentLengthGot == 0 && len(tc.TransferEncodingGot) == 0 &&
		!tc.IsUncompressed && tc.TLSGot == nil && tc.Time == 0
}

// iterate is the function to iterate the test cases
func (h *TestCases) iterate() <-chan TestCases {
	results := make(chan TestCases)
	go func() {
		defer close(results)

		current := h
		for current != nil {
			results <- *current
			current = current.Next
		}
	}()

	return results
}

// Add is the function to add test case
func (h *Instance) Add(testCase TestCase) error {
	var param = ""
	var body io.Reader

	// Set the query parameters
	if testCase.QueryParams != nil {
		param = testCase.QueryParams.(string)
	}

	// Set body parameters
	if testCase.BodyParams != nil {
		jsonBytes, err := json.Marshal(testCase.BodyParams)
		if err != nil {
			return err
		}

		body = bytes.NewBufferString(string(jsonBytes))
	}

	// Perform request
	url := h.Server.URL + testCase.EndPoint + param
	req, err := http.NewRequest(testCase.Method, url, body)
	if err != nil {
		return err
	}

	// Set body type
	if testCase.BodyType != "" {
		req.Header.Set("Content-Type", testCase.BodyType)
	}

	// Set headers
	if testCase.Header != nil {
		for key, value := range testCase.Header {
			req.Header.Set(key, value.(string))
		}
	}

	// Get response
	startTime := time.Now()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	if res == nil {
		return errors.New("response is nil")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	resGot, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Insert test case result
	return h.Cases.insert(
		TestCases{
			ResultGot:           string(resGot),
			StatusGot:           res.Status,
			StatusCodeGot:       res.StatusCode,
			ProtoGot:            res.Proto,
			ProtoMajorGot:       res.ProtoMajor,
			ProtoMinorGot:       res.ProtoMinor,
			ContentLengthGot:    res.ContentLength,
			TransferEncodingGot: res.TransferEncoding,
			IsUncompressed:      res.Uncompressed,
			TLSGot:              res.TLS,
			Time:                duration,
			Case:                testCase,
			Next:                nil,
		},
	)
}
