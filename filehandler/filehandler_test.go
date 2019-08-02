package filehandler

import (
	"gotest.tools/assert"
	"testing"
)

func TestBuildFilePathFromParts(t *testing.T) {
	testData := map[string][]string{
		"/a/b/c": {"/a", "/b", "/c"},
		"b/c/d":  {"b", "/c", "/d"},
		"/b/c/d": {"/b/", "/c/", "d"},
		"c/d/e":  {"c", "d", "e"},
		"d/e/f":  {"d/", "e/", "f"},
	}

	for key, value := range testData {
		got := BuildFilePathFromParts(value...)
		assert.Equal(t, got, key)
	}
}
