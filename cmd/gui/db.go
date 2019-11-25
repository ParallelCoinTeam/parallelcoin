package gui

import (
	"encoding/json"
	"fmt"
	scribble "github.com/nanobox-io/golang-scribble"
	"golang.org/x/text/unicode/norm"
	"unicode"
)

// Items

type DuOSitem struct {
	Enabled  bool        `json:"enabled"`
	Name     string      `json:"name"`
	Slug     string      `json:"slug"`
	Version  string      `json:"ver"`
	CompType string      `json:"comptype"`
	SubType  string      `json:"subtype"`
	Data     interface{} `json:"data"`
}
type DuOSitems struct {
	Slug  string              `json:"slug"`
	Items map[string]DuOSitem `json:"items"`
}

//  Vue App Model

type DuOScomp struct {
	IsApp       bool   `json:"isapp"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	Version     string `json:"ver"`
	Description string `json:"desc"`
	State       string `json:"state"`
	Image       string `json:"img"`
	URL         string `json:"url"`
	CompType    string `json:"comtype"`
	SubType     string `json:"subtype"`
	Js          string `json:"js"`
	Template    string `json:"template"`
	Css         string `json:"css"`
}

type DuOScomps []DuOScomp

type DuOSdb struct {
	DB     *scribble.Driver
	Folder string      `json:"folder"`
	Name   string      `json:"name"`
	Data   interface{} `json:"data"`
}

type Ddb interface {
	DbReadAllTypes()
	DbRead(folder, name string)
	DbReadAll(folder string) DuOSitems
	DbWrite(folder, name string, data interface{})
}

func (d *DuOSdb) DuoVueDbInit(dataDir string) {
	db, err := scribble.New(dataDir+"/gui", nil)
	if err != nil {
		fmt.Println("Error", err)
	}
	d.DB = db
}

var skip = []*unicode.RangeTable{
	unicode.Mark,
	unicode.Sk,
	unicode.Lm,
}

var safe = []*unicode.RangeTable{
	unicode.Letter,
	unicode.Number,
}

var _ Ddb = &DuOSdb{}

func (d *DuOSdb) DbReadAllTypes() {
	items := make(map[string]DuOSitems)
	types := []string{"assets", "config", "apps"}
	for _, t := range types {
		items[t] = d.DbReadAll(t)
	}
	d.Data = items
	fmt.Println("ooooooooooooooooooooooooooooodaaa", d.Data)

}
func (d *DuOSdb) DbReadTypeAll(f string) {
	d.Data = d.DbReadAll(f)
}

func (d *DuOSdb) DbReadAll(folder string) DuOSitems {
	itemsRaw, err := d.DB.ReadAll(folder)
	if err != nil {
		fmt.Println("Error", err)
	}
	items := make(map[string]DuOSitem)
	for _, bt := range itemsRaw {
		item := DuOSitem{}
		if err := json.Unmarshal([]byte(bt), &item); err != nil {
			fmt.Println("Error", err)
		}
		items[item.Slug] = item
	}
	return DuOSitems{
		Slug:  folder,
		Items: items,
	}
}

func (d *DuOSdb) DbReadAllComponents() map[string]DuOScomp {
	componentsRaw, err := d.DB.ReadAll("components")
	if err != nil {
		fmt.Println("Error", err)
	}
	components := make(map[string]DuOScomp)
	for _, componentRaw := range componentsRaw {
		component := DuOScomp{}
		if err := json.Unmarshal([]byte(componentRaw), &component); err != nil {
			fmt.Println("Error", err)
		}
		components[component.ID] = component
	}
	return components
}

func (d *DuOSdb) DbReadAddressBook() map[string]string {
	addressbookRaw, err := d.DB.ReadAll("addressbook")
	if err != nil {
		fmt.Println("Error", err)
	}
	addressbook := make(map[string]string)
	for _, addressRaw := range addressbookRaw {
		address := AddBook{}
		if err := json.Unmarshal([]byte(addressRaw), &address); err != nil {
			fmt.Println("Error", err)
		}
		addressbook[address.Address] = address.Label
	}
	return addressbook
}
func (d *DuOSdb) DbRead(folder, name string) {
	item := DuOSitem{}
	if err := d.DB.Read(folder, name, &item); err != nil {
		fmt.Println("Error", err)
	}
	d.Data = item
	fmt.Println("Daasdddddddaaa", item)
}
func (d *DuOSdb) DbWrite(folder, name string, data interface{}) {
	d.DB.Write(folder, name, data)
}

func slug(text string) string {
	buf := make([]rune, 0, len(text))
	dash := false
	for _, r := range norm.NFKD.String(text) {
		switch {
		case unicode.IsOneOf(safe, r):
			buf = append(buf, unicode.ToLower(r))
			dash = true
		case unicode.IsOneOf(skip, r):
		case dash:
			buf = append(buf, '-')
			dash = false
		}
	}
	if i := len(buf) - 1; i >= 0 && buf[i] == '-' {
		buf = buf[:i]
	}
	return string(buf)
}
