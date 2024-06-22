package apitest

import (
	"errors"
	"net/http"
	"net/http/httptest"
)

// Initialize creates a new instance of the API test
// with the given configuration and the need to exit
// the program after the test is done.
//
// Example:
//
//	config := &Config{
//		Level:        DefaultLevel,
//		IsNeedResult: true,
//	}
//	instance := apitest.Initialize(config, true)
func Initialize(config *Config, isNeedExit bool) *Instance {
	if err := ValidateConfig(config); err != nil {
		panic(err.Error())
		return nil
	}

	mux := http.NewServeMux()
	httptest.NewServer(mux)

	return &Instance{
		Server:           httptest.NewServer(mux),
		Mux:              mux,
		Cases:            &TestCases{},
		Config:           *config,
		TotalCases:       0,
		TotalFailedCases: 0,
		TotalPassedCases: 0,
		IsNeedExit:       isNeedExit,
	}
}

// Clean removes the instance of the API test.
// It closes the server and sets the mux, cases, and
// configuration to nil.
func (h *Instance) Clean() {
	h.Server.Close()
	h.Mux = nil
	h.Cases = nil
	h.Config = Config{}
}

// ValidateConfig validates the given configuration.
// It returns an error if the configuration is nil.
func ValidateConfig(config *Config) error {
	if config == nil {
		return errors.New("config is nil")
	}

	return nil
}
