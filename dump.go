package apitest

import (
	"fmt"
)

func dumpReadableHeader() string {
	return ""
}

func dumpReadableFooter() string {
	var output string
	output = fmt.Sprintf("\n%-23s : %d\n", "total test cases", 1)
	output += fmt.Sprintf("%-23s : %d\n", "total passed test cases", 1)
	output += fmt.Sprintf("%-23s : %d\n\n", "total failed test cases", 1)
	return output
}

func dumpReadableItem(result TestCases, config Config) string {
	var output string
	output = "\t\t{\n"
	output += "\t\t},\n"

	return output
}

func dumpJSONHeader() string {
	var output string
	output = "{\n"
	output += "\t\"test_cases\": [\n"
	return output
}

func dumpJSONFooter() string {
	var output string
	output = "\t],\n"
	output += fmt.Sprintf("\t\"total_test_cases\": %d,\n", 1)
	output += fmt.Sprintf("\t\"total_passed_test_cases\": %d,\n", 1)
	output += fmt.Sprintf("\t\"total_failed_test_cases\": %d,\n", 1)
	output += "}\n"
	return output
}

func dumpJSONItem(result TestCases, config Config) string {
	var output string
	// TODO

	return output
}

func dumpHeader(config Config) string {
	switch config.Format {
	case TableFormat:
		return dumpReadableHeader()
	case JSONFormat:
		return dumpJSONHeader()
	default:
		return dumpReadableHeader()
	}
}

func dumpItems(cases <-chan TestCases, config Config) string {
	var output = ""
	for result := range cases {
		switch config.Format {
		case TableFormat:
			output += dumpReadableItem(result, config)
		case JSONFormat:
			output += dumpJSONItem(result, config)
		default:
			output += dumpReadableItem(result, config)
		}
	}

	return output
}

func dumpFooter(config Config) string {
	switch config.Format {
	case TableFormat:
		return dumpReadableFooter()
	case JSONFormat:
		return dumpJSONFooter()
	default:
		return dumpReadableFooter()
	}
}

func performOutput(output string, config Config) error {
	// TODO

	return nil
}

// Dump is a method that writes the test cases to the destination based on the configuration.
// It takes an instance as an argument.
//
// Example:
//
//	instance := Instance
//	/* initialize instance .. */
//	instance.Cases.Dump(instance.Config)
func (h *TestCases) Dump(config Config) {
	var output string
	cases := h.iterate()
	output = dumpHeader(config)
	output += dumpItems(cases, config)
	output += dumpFooter(config)
	err := performOutput(output, config)
	if err != nil {
		panic(err)
	}
}
