package bnd

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"io/ioutil"
	"text/template"

	"os"
)

func filesLoop(b map[string]string) (fls string) {
	for f, bf := range b {
		fl := `"` + f + `":"` + bf + `",
`
		fls = fls + fl
	}
	return
}

// CompressFile reads the given file and converts it to a
// gzip compressed hex string
func CompressFile(filename string) (string, error) {
	data, err := ioutil.ReadFile("./assets/" + filename)
	if err != nil {
		return "", err
	}
	var byteBuffer bytes.Buffer
	writer := gzip.NewWriter(&byteBuffer)
	writer.Write(data)
	writer.Close()
	return hex.EncodeToString(byteBuffer.Bytes()), nil
}

// DecompressHexString decompresses the gzip/hex encoded data
func DecompressHexString(hexdata string) ([]byte, error) {
	data, err := hex.DecodeString(hexdata)
	if err != nil {
		panic(err)
	}
	datareader := bytes.NewReader(data)
	gzipReader, err := gzip.NewReader(datareader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()
	return ioutil.ReadAll(gzipReader)
}

func Bundle() map[string]string {
	b := make(map[string]string)
	for _, file := range files {
		zip, err := CompressFile(file)
		if err != nil {
		}
		b[file] = zip
	}
	var code = `package bnd
var FS = map[string]string{
` + filesLoop(b) + `}`

	file, _ := os.Create("./pkg/bnd/fs.go")
	defer file.Close()
	tmpl, _ := template.New("files").Parse(code)
	tmpl.Execute(file, "fs")
	return b
}

func DuOSsveBundler() sveBundle {
	fs := make(sveBundle)
	for f, fn := range Bundle() {
		unZip, err := DecompressHexString(fn)
		if err != nil {
		}
		fs[f] = unZip
	}
	return fs
}
