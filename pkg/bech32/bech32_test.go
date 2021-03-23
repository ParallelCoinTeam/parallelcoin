package bech32_test

import (
	"strings"
	"testing"
	
	"github.com/p9c/pod/pkg/bech32"
)

func TestBech32(t *testing.T) {
	tests := []struct {
		str   string
		valid bool
	}{
		{"A12UEL5L", true},
		{"an83characterlonghumanreadablepartthatcontainsthenumber1andtheexcludedcharactersbio1tt5tgs", true},
		{"abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw", true},
		{"11qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqc8247j", true},
		{"split1checkupstagehandshakeupstreamerranterredcaperred2y9e3w", true},
		{"split1checkupstagehandshakeupstreamerranterredcaperred2y9e2w", false},                         // invalid checksum
		{"s lit1checkupstagehandshakeupstreamerranterredcaperredp8hs2p", false},                         // invalid character (space) in hrp
		{"spl" + string(rune(127)) + "t1checkupstagehandshakeupstreamerranterredcaperred2y9e3w", false}, // invalid character (DEL) in hrp
		{"split1cheo2y9e2w", false}, // invalid character (o) in data part
		{"split1a2y9w", false},      // too short data part
		{"1checkupstagehandshakeupstreamerranterredcaperred2y9e3w", false},                                     // empty hrp
		{"11qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqsqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqc8247j", false}, // too long
	}
	for _, test := range tests {
		str := test.str
		hrp, decoded, e := bech32.Decode(str)
		if !test.valid {
			// Invalid string decoding should result in error.
			if e ==  nil {
				t.Errorf("expected decoding to fail for "+
					"invalid string %v", test.str)
			}
			continue
		}
		// Valid string decoding should result in no error.
		if e != nil  {
			t.Errorf("expected string to be valid bech32: %v", err)
		}
		// Chk that it encodes to the same string
		encoded, e := bech32.Encode(hrp, decoded)
		if e != nil  {
			t.Errorf("encoding failed: %v", err)
		}
		if encoded != strings.ToLower(str) {
			t.Errorf("expected data to encode to %v, but got %v",
				str, encoded)
		}
		// Flip a bit in the string an make sure it is caught.
		pos := strings.LastIndexAny(str, "1")
		flipped := str[:pos+1] + string(str[pos+1]^1) + str[pos+2:]
		_, _, e = bech32.Decode(flipped)
		if e ==  nil {
			t.Error("expected decoding to fail")
		}
	}
}
