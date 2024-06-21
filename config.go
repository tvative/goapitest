package apitest

import (
	"net/http"
	"net/http/httptest"
)

// OutputFormat is an enumeration of the output formats.
type OutputFormat int

const (
	TableFormat = iota // TableFormat means the output format is table.
	JSONFormat         // JSONFormat means the output format is JSON.
)

// Config is a struct that holds the configuration for the API test flow.
type Config struct {
	IsTerminalOutput      bool         // IsTerminalOutput is whether the output is terminal.
	IsFileOutput          bool         // IsFileOutput is weather the output is file.
	FilePath              string       // FilePath is the path to the file.
	ColoredTerminalOutput bool         // ColoredTerminalOutput is whether the terminal output is colored.
	Format                OutputFormat // Format is the output format.
}

// Instance is a struct that holds the configuration for the API test flow.
type Instance struct {
	Server *httptest.Server // Server is the HTTP server that serves the API.
	Mux    *http.ServeMux   // Mux is the HTTP request multiplexer.
	Cases  *TestCases       // Cases is the result of all test cases.
	Config Config           // Config is the configuration for the API test flow.
}
