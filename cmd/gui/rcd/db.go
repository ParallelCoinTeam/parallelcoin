package rcd

import (
	"encoding/json"
	"fmt"
	scribble "github.com/nanobox-io/golang-scribble"
	"golang.org/x/text/unicode/norm"
	"unicode"
)

// Items

type DuoUIitem struct {
	Enabled  bool        `json:"enabled"`
	Name     string      `json:"name"`
	Slug     string      `json:"slug"`
	Version  string      `json:"ver"`
	CompType string      `json:"comptype"`
	SubType  string      `json:"subtype"`
	Data     interface{} `json:"data"`
}
type DuoUIitems struct {
	Slug  string               `json:"slug"`
	Items map[string]DuoUIitem `json:"items"`
}

type DuoUIdb struct {
	DB     *scribble.Driver
	Folder string      `json:"folder"`
	Name   string      `json:"name"`
	Data   interface{} `json:"data"`
}

type Ddb interface {
	DbReadAllTypes()
	DbRead(folder, name string)
	DbReadAll(folder string) DuoUIitems
	DbWrite(folder, name string, data interface{})
}

func (d *DuoUIdb) DuoUIdbInit(dataDir string) {
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

var _ Ddb = &DuoUIdb{}

func (d *DuoUIdb) DbReadAllTypes() {
	items := make(map[string]DuoUIitems)
	types := []string{"assets", "config", "apps"}
	for _, t := range types {
		items[t] = d.DbReadAll(t)
	}
	d.Data = items
	fmt.Println("ooooooooooooooooooooooooooooodaaa", d.Data)

}
func (d *DuoUIdb) DbReadTypeAll(f string) {
	d.Data = d.DbReadAll(f)
}

func (d *DuoUIdb) DbReadAll(folder string) DuoUIitems {
	itemsRaw, err := d.DB.ReadAll(folder)
	if err != nil {
		fmt.Println("Error", err)
	}
	items := make(map[string]DuoUIitem)
	for _, bt := range itemsRaw {
		item := DuoUIitem{}
		if err := json.Unmarshal([]byte(bt), &item); err != nil {
			fmt.Println("Error", err)
		}
		items[item.Slug] = item
	}
	return DuoUIitems{
		Slug:  folder,
		Items: items,
	}
}

func (d *DuoUIdb) DbReadAddressBook() (addressbook map[string]string) {
	addressbook = make(map[string]string)
	if err := d.DB.Read("user", "addressbook", &addressbook); err != nil {
		fmt.Println("Error", err)
	}

	return addressbook
}
func (d *DuoUIdb) DbRead(folder, name string) {
	item := DuoUIitem{}
	if err := d.DB.Read(folder, name, &item); err != nil {
		fmt.Println("Error", err)
	}
	d.Data = item
}
func (d *DuoUIdb) DbWrite(folder, name string, data interface{}) {
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
