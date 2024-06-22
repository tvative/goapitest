package apitest

import (
	"time"
)

// TestCasePerformance is the performance of the test case.
type TestCasePerformance string

const (
	PerformanceWorst      TestCasePerformance = "\033[1;31mWorst\033[0m"      // PerformanceWorst means Low performance.
	PerformancePoor                           = "\033[1;35mPoor\033[0m"       // PerformancePoor means Medium performance.
	PerformanceAcceptable                     = "\033[1;33mAcceptable\033[0m" // PerformanceAcceptable means High performance.
	PerformanceGood                           = "\033[1;34mGood\033[0m"       // PerformanceGood means No performance.
	PerformanceBest                           = "\033[1;32mBest\033[0m"       // PerformanceBest means No performance.
)

type analyzedResult struct {
	IsPassed     bool
	PassedStatus string
	Performance  TestCasePerformance
	Type         string
}

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
