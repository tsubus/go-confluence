package export

// ProgressResponseForTest builds a ProgressResponse for external tests.
func ProgressResponseForTest(progress int, state, result string) ProgressResponse {
	return ProgressResponse{
		Progress: progress,
		State:    state,
		Result:   result,
	}
}

// EvaluateProgressForTest exposes evaluateProgress for external tests.
func EvaluateProgressForTest(pr ProgressResponse) (string, bool, error) {
	return evaluateProgress(pr)
}
