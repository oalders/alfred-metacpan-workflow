package wf_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/oalders/alfred-metacpan-workflow"
)

func TestSearchModule(t *testing.T) {
	result := wf.SearchModule("test")

	t.Log(result)

	if !strings.HasPrefix(result, "{") {
		t.Errorf("expected JSON output, got: %s", result)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Errorf("output is not valid JSON: %v", err)
	}

	if _, ok := parsed["items"]; !ok {
		t.Errorf("expected 'items' key in JSON output")
	}
}
