package cmd

import (
	"testing"
)

func Test_removeFromMap(t *testing.T) {
	testMap := map[string]string{"cat": "mouse"}
	removeFromMap(testMap, "cat")
	if _, ok := testMap["cat"]; ok != false {
		t.Errorf("Failed to remove 'cat' from map: %s", testMap)
	}
}
