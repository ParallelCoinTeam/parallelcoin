package btcjson_test

import (
	"encoding/json"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/btcjson"
)

// TestChainSvrWsResults ensures any results that have custom marshalling work as intended.
func TestChainSvrWsResults(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		result   interface{}
		expected string
	}{
		{
			name: "RescannedBlock",
			result: &btcjson.RescannedBlock{
				Hash:         "blockhash",
				Transactions: []string{"serializedtx"},
			},
			expected: `{"hash":"blockhash","transactions":["serializedtx"]}`,
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
