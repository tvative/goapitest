// GO API Testing Package
//
// Copyright (c) 2024 Tvative.
// All rights reserved.
//
// Source code and its usage are governed by
// the MIT license.

package apitest

import (
	"net/http"
	"net/http/httptest"
)

// OutputLevel is the type of output level.
type OutputLevel int

const (
	DefaultLevel = iota // DefaultLevel is the default output level.
	VerboseLevel        // VerboseLevel is the verbose output level.
)

// Config is the struct for configuration.
type Config struct {
	Level        OutputLevel // Level is the output level.
	IsNeedResult bool        // IsNeedResult is the flag to show the API result.
}

// Instance is the struct for the test instance
type Instance struct {
	Server           *httptest.Server // Server is the HTTP test server.
	Mux              *http.ServeMux   // Mux is the HTTP server multiplexer.
	Cases            *TestCases       // Cases is the test cases.
	Config           Config           // Config is the configuration.
	TotalCases       int64            // TotalCases is the total number of test cases.
	TotalFailedCases int64            // TotalFailedCases is the total number of failed test cases.
	TotalPassedCases int64            // TotalPassedCases is the total number of passed test cases.
	IsNeedExit       bool             // IsNeedExit is the flag to exit the test.
}
