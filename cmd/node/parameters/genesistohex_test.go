package parameters

import (
	"encoding/hex"
	"fmt"
	"testing"
)

var (
	mainnetGenesisHash, _     = hex.DecodeString(`000009f0fcbad3aac904d3660cfdcf238bf298cfe73adf1d39d14fc5c740ccc7`)
	mainnetGenesisBlock, _    = hex.DecodeString(`020000000000000000000000000000000000000000000000000000000000000000000000b79a9b6f31a9d7d25a1c4b0ec7a671dc56ce7663c380f2d2513a8e65e4ea43c8dcecc953ffff0f1e810201000101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff3a04ffff001d0104324e5954696d657320323031342d30372d3139202d2044656c6c20426567696e7320416363657074696e6720426974636f696effffffff0100e8764817000000434104e0d27172510c6806889740edafe6e63eb23fca32786fccfdb282bb2876a9f43b228245df057661ff943f6150716a20ea1851e8a7e9f54e620297664618438daeac00000000`)
	testnetGenesisHash, _     = hex.DecodeString(`00000e41ecbaa35ef91b0c2c22ed4d85fa12bbc87da2668fe17572695fb30cdf`)
	testnetGenesisBlock, _    = hex.DecodeString(`020000000000000000000000000000000000000000000000000000000000000000000000b79a9b6f31a9d7d25a1c4b0ec7a671dc56ce7663c380f2d2513a8e65e4ea43c884eac953ffff0f1e18df1a000101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff3a04ffff001d0104324e5954696d657320323031342d30372d3139202d2044656c6c20426567696e7320416363657074696e6720426974636f696effffffff0100e8764817000000434104e0d27172510c6806889740edafe6e63eb23fca32786fccfdb282bb2876a9f43b228245df057661ff943f6150716a20ea1851e8a7e9f54e620297664618438daeac00000000`)
	regtestnetGenesisHash, _  = hex.DecodeString(`69e9b79e220ea183dc2a52c825667e486bba65e2f64d237b578559ab60379181`)
	regtestnetGenesisBlock, _ = hex.DecodeString(`020000000000000000000000000000000000000000000000000000000000000000000000b79a9b6f31a9d7d25a1c4b0ec7a671dc56ce7663c380f2d2513a8e65e4ea43c8d4e5c953ffff7f20010000000101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff3a04ffff001d0104324e5954696d657320323031342d30372d3139202d2044656c6c20426567696e7320416363657074696e6720426974636f696effffffff0100e8764817000000434104e0d27172510c6806889740edafe6e63eb23fca32786fccfdb282bb2876a9f43b228245df057661ff943f6150716a20ea1851e8a7e9f54e620297664618438daeac00000000`)
)

func TestGenesisToHex(t *testing.T) {
	printByteAssignments("mainnetGenesisHash", *rev(mainnetGenesisHash))
	printByteAssignments("mainnetGenesisBlock", mainnetGenesisBlock)
	printByteAssignments("testnetGenesisHash", *rev(testnetGenesisHash))
	printByteAssignments("testnetGenesisBlock", testnetGenesisBlock)
	printByteAssignments("regtestnetGenesisHash", *rev(regtestnetGenesisHash))
	printByteAssignments("regtestnetGenesisBlock", regtestnetGenesisBlock)
}
func printByteAssignments(name string, in []byte) {
	fmt.Print(name, "=[]byte{\n")
	printGoHexes(in)
	fmt.Print("}\n")
}
func printGoHexes(in []byte) {
	fmt.Print("\t")
	for i := range in {
		if i%8 == 0 && i != 0 {
			fmt.Print("\n\t")
		}
		fmt.Printf("0x%02x, ", in[i])
	}
	fmt.Println()
}
func rev(in []byte) (out *[]byte) {
	o := make([]byte, len(in))
	out = &o
	for i := range in {
		(*out)[len(in)-i-1] = in[i]
	}
	return
}

// func hx(// 	in []byte) string {
// 	return hex.EncodeToString(in)
// }
// func split(// 	in []byte, pos int) (out []byte, piece []byte) {
// 	out = in[pos:]
// 	piece = in[:pos]
// 	return
// }
