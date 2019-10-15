package bnd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/p9c/pod/pkg/bundler/lib"
)

func Bundler(files []string) {

	var ignoreErrors = false

	bndFiles := lib.GetBundlerFiles(files, ignoreErrors)

	if len(bndFiles) == 0 {
		fmt.Println("No files found to process.")
		os.Exit(1)
	}
	referencedAssets, err := lib.GetReferencedAssets(bndFiles)
	if err != nil {
		log.Fatal(err)
	}

	targetFiles := []string{}

	for _, referencedAsset := range referencedAssets {
		packfileData, err := lib.GeneratePackFileString(referencedAsset, ignoreErrors)
		if err != nil {
			log.Fatal(err)
		}
		targetFile := filepath.Join(referencedAsset.BaseDir, referencedAsset.PackageName+"-bnd.go")
		targetFiles = append(targetFiles, targetFile)
		ioutil.WriteFile(targetFile, []byte(packfileData), 0644)
	}

	//var cmdargs []string

	//cmdargs = append(cmdargs, "build")
	//cmdargs = append(cmdargs, "-ldflags")
	//cmdargs = append(cmdargs, "-w -s")

	//cmd := exec.Command("go", cmdargs...)
	//stdoutStderr, err := cmd.CombinedOutput()
	//if err != nil {
	//	fmt.Printf("Error running command! %s\n", err.Error)
	//	fmt.Printf("From program: %s\n", stdoutStderr)
	//}

	// Remove target Files
	for _, filename := range targetFiles {
		err := os.Remove(filename)
		if err != nil {
			log.Fatal(err)
		}
	}

}
