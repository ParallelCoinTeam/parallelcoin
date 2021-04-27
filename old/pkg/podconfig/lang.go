package podconfig

import (
	"os"
)

func appLang(lang string) string {
	if lang == "" || lang == "en_US" {
		homeLang := os.Getenv("LANG")
		switch homeLang {
		case "sr_RS":
			lang = "rs"
		default:
			lang = "en"
		}
	}
	return lang
}

func Lang(lang string) string { return appLang(lang) }
