package main

import (
	"fmt"
	"github.com/p9c/pod/pkg/gui/wallet/translate/transifex"
	"github.com/p9c/pod/pkg/gui/wallet/translate/transifex/cli"
	"github.com/p9c/pod/pkg/gui/wallet/translate/transifex/config"
	"github.com/p9c/pod/pkg/gui/wallet/translate/transifex/format"
	"log"
	"path/filepath"
)

func main() {
	transifexCLI := cli.NewCLI()
	transifexApi := transifex.NewTransifexAPI(transifexCLI.ProjectSlug(), transifexCLI.Username(), transifexCLI.Password())
	rootDir := transifexCLI.RootDir()

	transifexApi.Debug = transifexCLI.Debug()

	var err error
	if err = transifexApi.ValidateConfiguration(); err != nil {
		log.Fatalf(err.Error())
	}

	var sourceLang string
	if sourceLang, err = transifexApi.SourceLanguage(); err != nil {
		log.Fatalf("Error loading the transifext project data.")
	}

	files, readFilesErr := config.ReadConfig(transifexCLI.ConfigFile(), rootDir, sourceLang)
	if readFilesErr != nil {
		fmt.Println(rootDir)
		log.Fatalf("Error reading reading language files: \n\n%s", readFilesErr)
	}

	existingResources := readExistingResources(transifexApi)

	doneChan := make(chan bool)
	goProcessNum := 0
	for _, file := range files {
		if _, has := existingResources[file.Slug]; has {
			goProcessNum++
			go downloadTranslations(rootDir, doneChan, sourceLang, file, transifexApi)
		}
	}

	for done := 0; done < goProcessNum; {
		<-doneChan

		done++
	}
}

func readExistingResources(transifexApi transifex.TransifexAPI) map[string]bool {
	resources, err := transifexApi.ListResources()
	if err != nil {
		log.Fatalf("Unable to load resources: %s", err)
	}
	existingResources := make(map[string]bool)
	for _, res := range resources {
		existingResources[res.Slug] = true
	}
	return existingResources
}

func downloadTranslations(rootDir string, doneChan chan bool, sourceLang string, file config.LocalizationFile, transifexApi transifex.TransifexAPI) {
	translations, err := transifexApi.DownloadTranslations(file.Slug)
	if err != nil {
		log.Fatalf("Failed to download translation files: %s", err)
	}
	i18Nformat, ok := format.Formats[file.I18nType]
	if !ok {
		log.Fatalf("No registered format %q", file.I18nType)
	}
	for _, path := range file.Translations {
		dir := filepath.Join(rootDir, filepath.Dir(path))
		for lang, translation := range translations {
			if err = i18Nformat().Write(dir, lang, sourceLang, file.Fname, translation, file.FileLocator); err != nil {
				log.Fatalf("Error writing out a translation: %s, %s\nError: %s\n\n Translation Data:\n%s", lang, file, err, translation)
			}
		}
		doneChan <- true
		break
	}
}
