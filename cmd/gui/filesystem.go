package gui

import (
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/shurcooL/vfsgen"
	"net/http"
)


func DuOSfileSystem(cx *conte.Xt) (err error) {
	// create filesystem
	var fs http.FileSystem = http.Dir("./pkg/gui/assets/filesystem")
	err = vfsgen.Generate(fs, vfsgen.Options{
		PackageName:  "guiFileSystem",
		BuildTags:    "dev",
		VariableName: "WalletGUI",
	})
	if err != nil {
		log.FATAL(err)
	}
	// add filesystem to bios struct
	cx.Gui.Fs = &fs
	// add bios struct to rcvar
	return
}