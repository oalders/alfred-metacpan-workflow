package wf_test

import (
	"testing"

	"github.com/oalders/alfred-metacpan-workflow"
)

func TestSearchModule(t *testing.T) {
	xmlStr := wf.SearchModule("test")

	t.Log(xmlStr)
}
