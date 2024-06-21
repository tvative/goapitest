package apitest

import (
	"errors"
	"net/http"
	"net/http/httptest"
)

// Initialize is a method that initializes the API test instance.
// It creates a new HTTP server and a new HTTP request multiplexer.
// It also initializes the test cases and the configuration.
// It takes a configuration as an argument.
//
// Example:
//
//	instance := Instance
//	config := Config{}
//	instance.Initialize(config)
func (h *Instance) Initialize(config *Config) {
	if err := ValidateConfig(config); err != nil {
		panic(err.Error())
		return
	}

	mux := http.NewServeMux()
	httptest.NewServer(mux)

	h.Server = httptest.NewServer(mux)
	h.Mux = mux
	h.Cases = nil
	h.Config = *config
}

// Clean is a method that cleans up the API test instance.
// It closes the HTTP server and sets the HTTP request multiplexer, test cases, and configuration to nil.
//
// Example:
//
//	instance := Instance
//	/* initialize instance ... */
//	instance.Clean()
func (h *Instance) Clean() {
	h.Server.Close()
	h.Mux = nil
	h.Cases = nil
	h.Config = Config{}
}

// ValidateConfig is a method that validates the configuration.
// It takes a configuration as an argument.
// It returns true if the configuration is valid.
//
// Example:
//
//	config := Config{}
//	valid := ValidateConfig(config)
func ValidateConfig(config *Config) error {
	if config == nil {
		return errors.New("config is nil")
	}

	if config.IsTerminalOutput && config.IsFileOutput {
		return errors.New("both terminal and file output are enabled")
	}

	if config.IsFileOutput && config.FilePath == "" {
		return errors.New("file path is empty")
	}

	switch config.Format {
	case TableFormat:
	case JSONFormat:
		break
	default:
		return errors.New("invalid output format")
	}

	if config.IsFileOutput && config.ColoredTerminalOutput {
		return errors.New("colored terminal output is enabled for file output")
	}

	return nil
}
