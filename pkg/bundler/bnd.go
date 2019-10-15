package bnd

import (
	"bytes"
	"github.com/p9c/pod/pkg/log"
	"text/template"

	"compress/gzip"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func (d *DuOSassets) DuOSassetsHandler() {
	for _, a := range *d {
		http.HandleFunc(a.Name, func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			headers["Content-Type"] = []string{a.ContentType}
			fmt.Fprint(w, DecompressHexString(a.Data))
		})
	}
}

func DuOSsveBundler() DuOSassets {
	a := DuOSassets{}
	for k, t := range Assets() {
		zip, err := CompressFile(Path(t))

		fmt.Println("paaaaaaaTTTTTT", Path(t))

		if err != nil {
		}
		t.Data = zip
		a[k] = t
	}
	var code = `package bnd
var FS = map[string]string{
` + filesLoop(a) + `}`

	file, _ := os.Create("./pkg/bundler/fs.go")
	defer file.Close()
	tmpl, _ := template.New("files").Parse(code)
	tmpl.Execute(file, "fs")
	return a
}

func filesLoop(b map[string]DuOSasset) (fls string) {
	for f, bf := range b {
		fl := `"` + f + `":"` + bf.Data + `",
`
		fls = fls + fl
	}
	return
}

func Path(a DuOSasset) (path string) {
	if a.Sub != "" {
		path = a.Sub + "/" + a.Name
	} else {
		path = a.Name
	}
	return
}

// CompressFile reads the given file and converts it to a
// gzip compressed hex string
func CompressFile(p string) (string, error) {
	data, err := ioutil.ReadFile("./assets/" + p)
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
func DecompressHexString(hexdata string) string {
	data, err := hex.DecodeString(hexdata)
	if err != nil {
		log.ERROR(err)
	}
	datareader := bytes.NewReader(data)
	gzipReader, err := gzip.NewReader(datareader)
	if err != nil {
		log.ERROR(err)
	}
	defer gzipReader.Close()
	asset, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		log.ERROR(err)
	}
	return string(asset)
}
