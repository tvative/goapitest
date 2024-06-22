<div align="center">
  <h1>Go API Test</h1>
  <p>A Simple and Minimalist API Testing Library for Go Language</p>
</div>

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/github.com/tvative/goapitest)

`goapitest` is a minimalist package to test RESTful APIs. It is designed to be simple and easy to use, with a focus on
simplicity and minimalism.

## Usage

### Installation

To use `goapitest` in your Go project, you need to install it using Go modules. You can add it to your project with the
following command:

```bash
go get -u github.com/tvative/goapitest
```

### Basic Usage

To use `goapitest`, you need to create a new instance of the `API` struct. The complete code is as follows:

```go
package main

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
			Details:        "Happy path test case 01",
			EndPoint:       "/hello",
			Method:         http.MethodGet,
			Header:         nil,
			BodyType:       "",
			QueryParams:    nil,
			BodyParams:     nil,
			ExpectedResult: "",
			ExpectedStatus: http.StatusOK,
			IsIgnored:      false,
		}
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
	// Content of the handler..
}
```

### Example

For a minimal example, see the [test/main.go](test/main.go) file. That file contains a simple example of how to use the
package.

## License

This project is licensed under the MIT License; see the [LICENSE](LICENSE) file for details.
