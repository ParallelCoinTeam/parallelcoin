package pkger

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/p9c/pod/pkg/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var code = `
package assets

var PkgBundle  = map[string][]byte {
{{range $key, $val := .}}
     "{{$key}}":{{$val}},
{{end}}
}
`



func Pkger(dirname string) {
	pkgBundle := make(map[string][]byte)

	err := filepath.Walk(dirname,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				fmt.Println(path)
				asset, err := ioutil.ReadFile(path) // b has type []byte
				if err != nil {
					log.ERROR("error building packages:", err)
				}
				var assetZip bytes.Buffer

				w := gzip.NewWriter(&assetZip)
				w.Write(asset)
				w.Close()

				//r, err := gzip.NewReader(&b)
				//io.Copy(os.Stdout, r)
				//r.Close()
				pkgBundle[strings.TrimLeft(path, "pkg/gui/assets/filesystem/")] = assetZip.Bytes()

			} else {

			}
			return nil
		})
	if err != nil {
		log.ERROR("error building packages:", err)
	}

	//flag.Parse()
	file, _ := os.Create("./pkg/gui/assets/pkgbundle.go")
	defer file.Close()

	tmpl, _ := template.New("").Parse(code)
	tmpl.Execute(file, pkgBundle)

}

//func Pkger(dirname string) {
//
//
//	rootDir, err := subDir(dirname)
//	if err != nil {
//		log.ERROR("error building packages:", err)
//	}
//	for _, file := range rootDir {
//		if file.IsDir(){
//			fmt.Println("isDir --> " + file.Name())
//			subFolder, err := subDir(dirname)
//			if err != nil {
//				log.ERROR("error building packages:", err)
//			}
//		}else {
//			fmt.Println("isFile --> " + file.Name())
//		}
//	}
//}

//
//
//func subDir(d string) ([]os.FileInfo, error){
//	f, err := os.Open(d)
//	if err != nil {
//		log.FATAL("error opening folder:", err)
//	}
//	files, err := f.Readdir(-1)
//	f.Close()
//	if err != nil {
//		log.FATAL("error building packages:", err)
//	}
//return files, err
//}
//
