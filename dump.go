// GO API Testing Package
//
// Copyright (c) 2024 Tvative.
// All rights reserved.
//
// Source code and its usage are governed by
// the MIT license.

package apitest

import (
	"fmt"
	"os"
	"strings"
)

// dumpReadableHeader is the function to dump the readable headers
func dumpReadableHeader() string {
	return fmt.Sprintf("Running API testing...\n\n")
}

// dumpReadableFooter is the function to dump the readable footer
func dumpReadableFooter(instance *Instance) string {
	var output string
	output = fmt.Sprintf("%-23s : %d\n", "total test cases", instance.TotalCases)
	output += fmt.Sprintf("%-23s : %d/%d\n", "total passed test cases", instance.TotalPassedCases, instance.TotalCases)
	output += fmt.Sprintf("%-23s : %d/%d\n", "total failed test cases", instance.TotalFailedCases, instance.TotalCases)
	return output
}

// dumpReadableItem is the function to dump the readable item
func dumpReadableItem(tc TestCases, config Config, result analyzedResult) string {
	var output string
	if result.IsPassed {
		output = fmt.Sprintf("Test case %-7s with %s performance [ %s ]\n", result.PassedStatus, result.Performance, tc.Case.ID)
	} else {
		output = fmt.Sprintf("Test case %-7s [ %s ]\n", result.PassedStatus, tc.Case.ID)
	}

	output += fmt.Sprintf("  ├─ %-15s : %s\n", "Details", tc.Case.Details)
	output += fmt.Sprintf("  ├─ %-15s : %s\n", "Endpoint", tc.Case.EndPoint)
	output += fmt.Sprintf("  ├─ %-15s : %s\n", "Type", result.Type)

	if result.IsPassed && !config.IsNeedResult {
		output += fmt.Sprintf("  └─ %-15s : %s\n", "Time", tc.Time)
	}

	if result.IsPassed && config.IsNeedResult {
		output += fmt.Sprintf("  ├─ %-15s : %s\n", "Time", tc.Time)
		output += fmt.Sprintf("  └─ %-15s : %s\n", "Result", strings.TrimSuffix(tc.ResultGot, "\n"))
	}

	if !result.IsPassed {
		output += fmt.Sprintf("  ├─ %-15s : %s\n", "Time", tc.Time)
		output += fmt.Sprintf("  ├─ %-15s : %s\n", "Performance", result.Performance)
		output += fmt.Sprintf("  ├─ %-15s : Expected %d, but got \033[1;31m%d\033[0m (%s)\n", "Status", tc.Case.ExpectedStatus, tc.StatusCodeGot, tc.StatusGot)

		if config.Level == VerboseLevel {
			output += fmt.Sprintf("  ├─ %-15s : %s (%d, %d)\n", "Protocol", tc.ProtoGot, tc.ProtoMajorGot, tc.ProtoMinorGot)
			output += fmt.Sprintf("  ├─ %-15s : %d\n", "Content Length", tc.ContentLengthGot)
		}

		output += fmt.Sprintf("  └─ %-15s : %s\n", "Result", strings.TrimSuffix(tc.ResultGot, "\n"))
	}

	return output + "\n"
}

// dumpHeader is the function to dump the header
func dumpHeader() string {
	return dumpReadableHeader()
}

// dumpItems is the function to dump the items
func dumpItems(cases <-chan TestCases, instance *Instance) string {
	var output = ""
	for tc := range cases {
		if isEmptyNode(tc) {
			continue
		}

		result := instance.analyze(tc)
		output += dumpReadableItem(tc, instance.Config, result)
	}

	return output
}

// dumpFooter is the function to dump the footer
func dumpFooter(instance *Instance) string {
	return dumpReadableFooter(instance)
}

// performOutput is the function to perform the output
func performOutput(output string, instance *Instance) error {
	fmt.Print(output)

	if instance.IsNeedExit {
		if instance.TotalFailedCases > 0 {
			os.Exit(1)
		}
	}

	return nil
}

// Dump is the function to dump the test cases
// in readable format
// It will print the test cases to the console
//
// Example:
//
//	instance.Cases.Dump(instance)
func (h *TestCases) Dump(instance *Instance) {
	var output string
	cases := h.iterate()
	output = dumpHeader()
	output += dumpItems(cases, instance)
	output += dumpFooter(instance)
	err := performOutput(output, instance)
	if err != nil {
		panic(err)
	}
}
