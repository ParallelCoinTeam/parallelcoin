package btcjson_test

import (
	"encoding/json"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/btcjson"
)

// TestPodExtCustomResults ensures any results that have custom marshalling work as intended and unmarshal code of
// results are as expected.
func TestPodExtCustomResults(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		result   interface{}
		expected string
	}{
		{
			name: "versionresult",
			result: &btcjson.VersionResult{
				VersionString: "1.0.0",
				Major:         1,
				Minor:         0,
				Patch:         0,
				Prerelease:    "pr",
				BuildMetadata: "bm",
			},
			expected: `{"versionstring":"1.0.0","major":1,"minor":0,"patch":0,"prerelease":"pr","buildmetadata":"bm"}`,
		},
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		marshalled, e := json.Marshal(test.result)
		if e != nil {
			t.Errorf("Test #%d (%s) unexpected error: %v", i,
				test.name, e,
			)
			continue
		}
		if string(marshalled) != test.expected {
			t.Errorf("Test #%d (%s) unexpected marhsalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.expected,
			)
			continue
		}
	}
}
