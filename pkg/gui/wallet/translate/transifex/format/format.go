package format

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// A Format is responsible for cleaning up the translation file.  For example:
// * A json file with empty string key should be deleted
// * A json file with empty values should be converted to a string with a single space
// * An xml file with nested tags should be flattened
type Format interface {
	// initialize the format with any extra paramaters available in the configuration
	Init(map[string]interface{})
	// Returns the extension of the format files
	Ext() string
	// takes the raw data read from file
	// returns:
	// * cleaned bytes
	// * new i18nType - it is possible that the format decided that in order to make the format useable it needed to convert it to a new type
	// * an error or nil if no errors occurred
	Clean([]byte) ([]byte, string, error)
	// Write a new translation to the correct translation file
	// * rootDir - path to the root of the translation files directory tree
	// * langCode - the language code of the translation
	// * srcLang - the source language
	// * translation - the translation text
	// * fileLocator - strategy for creating the path to the translation file
	Write(rootDir, langCode, srcLang, filename, translation string, fileLocator FileLocator) error
}

// factory methods for creating a format object
var Formats = map[string]func() Format{
	"KEYVALUEJSON":     func() Format { return new(KeyValueJson) },
	"FLATTENXMLTOJSON": func() Format { return new(FlattenXmlToJson) }}

// Strategy for creating the path to a translation file
type FileLocator interface {
	// path - the path to the directory containing the translation files (or root of the tree for looking up the files)
	// lang - the language of the translation file
	// name - the name of the file (or base name)
	// ext - the extension of the file
	Find(path, lang, name, ext string) string
	List(path, name, ext string) (map[string]string, error)
}

var identityMapper = func(lang string, reverse bool) string { return lang }

var langCodeMapper2To3 = func(lang string, reverse bool) string {
	var mapping map[string]string
	if reverse {
		mapping = threeToTwoLetterIsoCode
	} else {
		mapping = twoToThreeLetterIsoCode
	}
	code, has := mapping[lang]
	if !has {
		panic("There is no known language code mapper for: " + lang)
	}

	return code
}
var FileLocators = map[string]FileLocator{
	"LANG-NAME":        LangNameLocator{identityMapper},
	"3-CHAR-LANG-NAME": LangNameLocator{langCodeMapper2To3},
	"3-CHAR-LOC-DIR":   LocDirLocator{langCodeMapper2To3},
	"LOC-DIR":          LocDirLocator{identityMapper}}

// All translation files in same directory and have pattern: lang-name.ext
type LangNameLocator struct {
	// map the 2 letter language code to the required code for lookup
	// first param is the letter to map
	// second param is the direction of the mapping (ie 2letter -> 3letter or vice-versa)
	langCodeMapper func(string, bool) string
}

func (l LangNameLocator) Find(path, lang, name, ext string) string {
	flang := l.langCodeMapper(lang, false)
	fileName := fmt.Sprintf("%s-%s.%s", flang, name, ext)
	return filepath.Join(path, fileName)
}

func (l LangNameLocator) List(path, name, ext string) (map[string]string, error) {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	filename := fmt.Sprintf("-%s.%s", name, ext)

	candidates, readErr := ioutil.ReadDir(path)

	if readErr != nil {
		return nil, readErr
	}

	translationFiles := map[string]string{}

	for _, f := range candidates {
		fname := filepath.Join(path, f.Name())

		if strings.HasSuffix(fname, filename) {
			lang := strings.Split(filepath.Base(fname), "-")[0]
			lang = l.langCodeMapper(lang, true)
			translationFiles[lang] = fname
		}
	}

	return translationFiles, nil
}

// Locate translation file with a directory structure as follows:
// <root>
// -- en
//    -- name.<ext>
// -- fr
//    -- name.<ext>
type LocDirLocator struct {
	// map the 2 letter language code to the required code for lookup
	// first param is the letter to map
	// second param is the direction of the mapping (ie 2letter -> 3letter or vice-versa)
	langCodeMapper func(string, bool) string
}

func (l LocDirLocator) Find(path, lang, name, ext string) string {
	flang := l.langCodeMapper(lang, false)
	fileName := fmt.Sprintf("%s/%s.%s", flang, name, ext)
	return filepath.Join(path, fileName)
}

func (l LocDirLocator) List(path, name, ext string) (map[string]string, error) {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	filename := fmt.Sprintf("%s.%s", name, ext)

	locDirs, readErr := ioutil.ReadDir(path)

	if readErr != nil {
		return nil, readErr
	}

	translationFiles := map[string]string{}

	for _, locDir := range locDirs {
		loc := locDir.Name()
		fname := filepath.Join(path, loc, filename)
		_, err := os.Stat(fname)

		if err == nil {
			loc = l.langCodeMapper(loc, true)
			translationFiles[loc] = fname
		}
	}

	return translationFiles, nil
}

var threeToTwoLetterIsoCode map[string]string

func init() {
	threeToTwoLetterIsoCode = make(map[string]string)
	for key, value := range twoToThreeLetterIsoCode {
		threeToTwoLetterIsoCode[value] = key
	}
}

var twoToThreeLetterIsoCode = map[string]string{
	"aa": "aar",
	"ab": "abk",
	"af": "afr",
	"ak": "aka",
	"sq": "alb",
	"am": "amh",
	"ar": "ara",
	"an": "arg",
	"hy": "arm",
	"as": "asm",
	"av": "ava",
	"ae": "ave",
	"ay": "aym",
	"az": "aze",
	"ba": "bak",
	"bm": "bam",
	"eu": "baq",
	"be": "bel",
	"bn": "ben",
	"bh": "bih",
	"bi": "bis",
	"bs": "bos",
	"br": "bre",
	"bg": "bul",
	"my": "bur",
	"ca": "cat",
	"ch": "cha",
	"ce": "che",
	"zh": "chi",
	"cu": "chu",
	"cv": "chv",
	"kw": "cor",
	"co": "cos",
	"cr": "cre",
	"cs": "cze",
	"da": "dan",
	"dv": "div",
	"nl": "dut",
	"dz": "dzo",
	"en": "eng",
	"eo": "epo",
	"et": "est",
	"ee": "ewe",
	"fo": "fao",
	"fj": "fij",
	"fi": "fin",
	"fr": "fre",
	"fy": "fry",
	"ff": "ful",
	"ka": "geo",
	"de": "ger",
	"gd": "gla",
	"ga": "gle",
	"gl": "glg",
	"gv": "glv",
	"el": "gre",
	"gn": "grn",
	"gu": "guj",
	"ht": "hat",
	"ha": "hau",
	"he": "heb",
	"hz": "her",
	"hi": "hin",
	"ho": "hmo",
	"hu": "hun",
	"ig": "ibo",
	"is": "ice",
	"io": "ido",
	"ii": "iii",
	"iu": "iku",
	"ie": "ile",
	"ia": "ina",
	"id": "ind",
	"ik": "ipk",
	"it": "ita",
	"jv": "jav",
	"ja": "jpn",
	"kl": "kal",
	"kn": "kan",
	"ks": "kas",
	"kr": "kau",
	"kk": "kaz",
	"km": "khm",
	"ki": "kik",
	"rw": "kin",
	"ky": "kir",
	"kv": "kom",
	"kg": "kon",
	"ko": "kor",
	"kj": "kua",
	"ku": "kur",
	"lo": "lao",
	"la": "lat",
	"lv": "lav",
	"li": "lim",
	"ln": "lin",
	"lt": "lit",
	"lb": "ltz",
	"lu": "lub",
	"lg": "lug",
	"mk": "mac",
	"mh": "mah",
	"mi": "mao",
	"mr": "mar",
	"ms": "may",
	"mg": "mlg",
	"mt": "mlt",
	"ml": "mol",
	"mn": "mon",
	"na": "nau",
	"nv": "nav",
	"nr": "nbl",
	"nd": "nde",
	"ng": "ndo",
	"ne": "nep",
	"nn": "nno",
	"nb": "nob",
	"no": "nor",
	"ny": "nya",
	"oc": "oci",
	"oj": "oji",
	"or": "ori",
	"om": "orm",
	"os": "oss",
	"pa": "pan",
	"fa": "per",
	"pi": "pli",
	"pl": "pol",
	"pt": "por",
	"ps": "pus",
	"qu": "que",
	"rm": "roh",
	"ro": "rum",
	"rn": "run",
	"ru": "rus",
	"sg": "sag",
	"sa": "san",
	"sr": "srp",
	"hr": "hrv",
	"si": "sin",
	"sk": "slo",
	"sl": "slv",
	"se": "sme",
	"sm": "smo",
	"sn": "sna",
	"sd": "snd",
	"so": "som",
	"st": "sot",
	"es": "spa",
	"sc": "srd",
	"ss": "ssw",
	"su": "sun",
	"sw": "swa",
	"sv": "swe",
	"ty": "tah",
	"ta": "tam",
	"tt": "tat",
	"te": "tel",
	"tg": "tgk",
	"tl": "tgl",
	"th": "tha",
	"bo": "tib",
	"ti": "tir",
	"to": "ton",
	"tn": "tsn",
	"ts": "tso",
	"tk": "tuk",
	"tr": "tur",
	"tw": "twi",
	"ug": "uig",
	"uk": "ukr",
	"ur": "urd",
	"uz": "uzb",
	"ve": "ven",
	"vi": "vie",
	"vo": "vol",
	"cy": "wel",
	"wa": "wln",
	"wo": "wol",
	"xh": "xho",
	"yi": "yid",
	"yo": "yor",
	"za": "zha",
	"zu": "zul"}
