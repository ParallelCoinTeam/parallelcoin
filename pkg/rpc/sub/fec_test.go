package sub

import (
	"encoding/binary"
	"encoding/hex"
	"testing"
)

var (
	testDataAligned   = []byte("123456789123456789123456789123456789123456789123456789123456789123456789123456789")
	testDataUnaligned = []byte("1234567891234567891234567891234567891234")
	expectedAligned   = "510031323334353637383931323334353637383931323334353637383931323334353637383931323334353637383931323334353637383931323334353637383931323334353637383931323334353637383900000000000000"
	expectedUnaligned = "280031323334353637383931323334353637383931323334353637383931323334353637383931323334000000"
)

func TestPadData(	t *testing.T) {
	actualAligned := hex.EncodeToString(padData(testDataAligned))
	actualUnaligned := hex.EncodeToString(padData(testDataUnaligned))
	if actualAligned != expectedAligned {
		t.Fatalf("Padding did not produce expected result:\ngot      '%s'\nexpected '%s'",
			actualAligned, expectedAligned)
	}
	if actualUnaligned != expectedUnaligned {
		t.Fatalf("Padding did not produce expected result:\ngot      '%s'\nexpected '%s'",
			actualUnaligned, expectedUnaligned)
	}
}
func TestFECCodec(	t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("Recovered in f", r)
		}
	}()
	chunks := rsEncode(testDataAligned)
	// Deface one of the pieces
	chunks[4][3] = ^chunks[4][3]
	// Here we only need 3 packets
	_, err := rsDecode(chunks[4:7])
	if err != nil {
		panic(err)
	}
	// Requires one more across the punctured chunk to recover. This would not normally happen as the checksums would usually filter out incorrect chunks.
	var data []byte
	data, err = rsDecode(chunks[3:6])
	if err != nil {
		panic(err)
	}
	dataLen := binary.LittleEndian.Uint16(data)
	result := data[2 : dataLen+2]
	dataString := hex.EncodeToString(data[2 : dataLen+2])
	resultString := hex.EncodeToString(result)
	if dataString != resultString {
		t.Fatalf("FEC encode/decode failed:\ngot      '%s'\nexpected '%s'", dataString, resultString)
	}
}
