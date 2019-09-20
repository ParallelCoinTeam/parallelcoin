// Package graph turns hexadecimal strings into visual representation of the
// magnitude of digits
package graf

var table = map[string][]string{
	"0": {" ", " "},
	"1": {" ", "▁"},
	"2": {" ", "▂"},
	"3": {" ", "▃"},
	"4": {" ", "▄"},
	"5": {" ", "▅"},
	"6": {" ", "▆"},
	"7": {" ", "▇"},
	"8": {" ", "█"},
	"9": {"▁", "█"},
	"a": {"▂", "█"},
	"b": {"▃", "█"},
	"c": {"▄", "█"},
	"d": {"▅", "█"},
	"e": {"▆", "█"},
	"f": {"▇", "█"},
}

func Convert(hexString string) (out string) {
	line1, line2 := "", ""
	for i := range hexString {
		line1 += table[string(hexString[i])][0]
		line2 += table[string(hexString[i])][1]
	}
	out = "\n" + line1 + "\n" + line2
	return
}
