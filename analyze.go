// GO API Testing Package
//
// Copyright (c) 2024 Tvative.
// All rights reserved.
//
// Source code and its usage are governed by
// the MIT license.

package apitest

import (
	"time"
)

// TestCasePerformance is the type of test case performance.
type TestCasePerformance string

const (
	PerformanceWorst      TestCasePerformance = "\033[1;31mWorst\033[0m"      // PerformanceWorst is the worst performance and the color is red.
	PerformancePoor                           = "\033[1;35mPoor\033[0m"       // PerformancePoor is the poor performance and the color is magenta.
	PerformanceAcceptable                     = "\033[1;33mAcceptable\033[0m" // PerformanceAcceptable is the acceptable performance and the color is yellow.
	PerformanceGood                           = "\033[1;34mGood\033[0m"       // PerformanceGood is the good performance and the color is blue.
	PerformanceBest                           = "\033[1;32mBest\033[0m"       // PerformanceBest is the best performance and the color is green.
)

// analyzedResult is the struct for analyzed result
type analyzedResult struct {
	IsPassed     bool                // IsPassed is the flag to check if the test case is passed.
	PassedStatus string              // PassedStatus is the status of the test case.
	Performance  TestCasePerformance // Performance is the performance of the test case.
	Type         string              // Type is the type of the test case.
}

// analyze is the function to analyze the test case.
func (h *Instance) analyze(tc TestCases) analyzedResult {
	var result analyzedResult
	if tc.Case.ExpectedStatus != tc.StatusCodeGot {
		result.IsPassed = false
		result.PassedStatus = "\u001B[1;31mFailed\u001B[0m"
	} else {
		result.IsPassed = true
		result.PassedStatus = "\u001B[1;32mPassed\u001B[0m"
	}

	if tc.Case.ExpectedResult != "" {
		if tc.Case.ExpectedResult != tc.ResultGot {
			result.IsPassed = false
			result.PassedStatus = "\u001B[1;31mFailed\u001B[0m"
		}
	}

	result.Performance = getPerformance(tc.Time)
	result.Type = getType(tc.Case.Type)

	h.TotalCases++
	if result.IsPassed {
		h.TotalPassedCases++
	} else {
		h.TotalFailedCases++
	}

	return result
}

// getPerformance is the function to get the performance.
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

// getType is the function to get the type.
func getType(caseType TestCaseType) string {
	switch caseType {
	case HappyPath:
		return "Happy path"
	case EdgeCase:
		return "Edge case"
	case NegativeCase:
		return "Negative case"
	case BoundaryCase:
		return "Boundary case"
	case CornerCase:
		return "Corner case"
	case StressCase:
		return "Stress case"
	case SmokeCase:
		return "Smoke case"
	case RegressionCase:
		return "Regression case"
	case IntegrationCase:
		return "Integration case"
	default:
		return "Unknown"
	}
}
