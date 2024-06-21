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

// TestCaseType is the type of test case.
type TestCaseType int

const (
	HappyPath       TestCaseType = iota // HappyPath means Happy path test case.
	EdgeCase                            // EdgeCase means Edge case test case.
	NegativeCase                        // NegativeCase means Negative case test case.
	BoundaryCase                        // BoundaryCase means Boundary case test case.
	CornerCase                          // CornerCase means Corner case test case.
	StressCase                          // StressCase means Stress case test case.
	SmokeCase                           // SmokeCase means Smoke case test case.
	RegressionCase                      // RegressionCase means Regression case test case.
	IntegrationCase                     // IntegrationCase means Integration case test case.
)

// TestCasePerformance is the performance of the test case.
type TestCasePerformance string

const (
	PerformanceWorst      TestCasePerformance = "Worst"      // PerformanceWorst means Low performance.
	PerformancePoor                           = "Poor"       // PerformancePoor means Medium performance.
	PerformanceAcceptable                     = "Acceptable" // PerformanceAcceptable means High performance.
	PerformanceGood                           = "Good"       // PerformanceGood means No performance.
	PerformanceBest                           = "Best"       // PerformanceBest means No performance.
)

// TestCase is a single test case
type TestCase struct {
	ID             string         // ID is the unique identifier for the test case.
	Type           TestCaseType   // Type is the type of the test case.
	Details        string         // Details is the description of the test case.
	EndPoint       string         // EndPoint is the endpoint to test.
	Method         string         // Method is the HTTP method to use.
	Header         map[string]any // Header is the headers to send.
	BodyType       string         // BodyType is the body type to send.
	QueryParams    any            // QueryParams is the query parameters to send.
	BodyParams     any            // BodyParams is the body parameters to send.
	ExpectedResult string         // ExpectedResult is the expected result.
	ExpectedStatus int            // ExpectedStatus is the expected HTTP status code.
	IsIgnored      bool           // IsIgnored is whether the test case is ignored.
}

// TestCases is the result of a test case.
type TestCases struct {
	ResultGot           string               // ResultGot is the result got from the test.
	StatusGot           string               // StatusGot is the status got from the test.
	StatusCodeGot       int                  // ResultGot is the status code got from the test.
	ProtoGot            string               // ProtoGot is the protocol got from the test.
	ProtoMajorGot       int                  // ProtoMajorGot is the major protocol version got from the test.
	ProtoMinorGot       int                  // ProtoMinorGot is the minor protocol version got from the test.
	ContentLengthGot    int64                // ContentLengthGot is the content length got from the test.
	TransferEncodingGot []string             // TransferEncodingGot is the transfer encoding got from the test.
	IsUncompressed      bool                 // IsUncompressed is whether the response is uncompressed.
	TLSGot              *tls.ConnectionState // TLSGot is the TLS connection state got from the test.
	Time                time.Duration        // Time is the time taken to run the test case.
	Performance         TestCasePerformance  // Performance is the performance of the test case.
	Case                TestCase             // Case is the test case.
	Next                *TestCases           // Next is the next test case result.
}

func getPerformance(resTime time.Duration) TestCasePerformance {
	if resTime < 100*time.Millisecond {
		return PerformanceBest
	}

	if resTime >= 100*time.Millisecond && resTime < 300*time.Millisecond {
		return PerformanceGood
	}

	if resTime >= 300*time.Millisecond && resTime < 1000*time.Millisecond {
		return PerformanceAcceptable
	}

	if resTime >= 1000*time.Millisecond && resTime < 2000*time.Millisecond {
		return PerformancePoor
	}

	return PerformanceWorst
}

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

// iterate returns a channel that will receive test case results one by one.
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

// Add adds a test case to the configuration.
//
// Example:
//
//	instance := Instance
//	/* initialize instance .. */
//	testCase := TestCase{}
//	instance.Add(testCase)
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
			Performance:         getPerformance(duration),
			Case:                testCase,
			Next:                nil,
		},
	)
}
