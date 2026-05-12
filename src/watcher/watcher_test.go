package watcher

import "testing"

func TestExtractAgentNameFromHandoffFilename(t *testing.T) {
	tests := map[string]string{
		"SYSTEM_INIT_TO_CEO_PHASE_KICKOFF_120000.json": "CEO",
		"CEO_TO_DEV_BACKEND_TASK_001.json":             "DEV_BACKEND",
		"review_to_CODE_REVIEWER.md":                   "CODE_REVIEWER",
		"no_recipient.json":                            "",
	}

	for filename, want := range tests {
		if got := ExtractAgentName(filename); got != want {
			t.Fatalf("ExtractAgentName(%q) = %q, want %q", filename, got, want)
		}
	}
}
