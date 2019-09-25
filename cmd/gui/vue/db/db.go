//+build !nogui
// +build !headless

package db

import (
	"encoding/json"
	"fmt"
	"unicode"

	"github.com/p9c/pod/cmd/gui/vue/mod"

	"golang.org/x/text/unicode/norm"
)

var skip = []*unicode.RangeTable{
	unicode.Mark,
	unicode.Sk,
	unicode.Lm,
}

var safe = []*unicode.RangeTable{
	unicode.Letter,
	unicode.Number,
}

var _ DVdb = &DuoVUEdb{}

func (d *DuoVUEdb) DbReadAllTypes() {
	items := make(map[string]mod.DuoGuiItems)
	types := []string{"assets", "config", "apps"}
	for _, t := range types {
		items[t] = d.DbReadAll(t)
	}
	d.Data = items
	fmt.Println("ooooooooooooooooooooooooooooodaaa", d.Data)

}
func (d *DuoVUEdb) DbReadTypeAll(f string) {
	d.Data = d.DbReadAll(f)
}

func (d *DuoVUEdb) DbReadAll(folder string) mod.DuoGuiItems {
	itemsRaw, err := d.DB.ReadAll(folder)
	if err != nil {
		fmt.Println("Error", err)
	}
	items := make(map[string]mod.DuoGuiItem)
	for _, bt := range itemsRaw {
		item := mod.DuoGuiItem{}
		if err := json.Unmarshal([]byte(bt), &item); err != nil {
			fmt.Println("Error", err)
		}
		items[item.Slug] = item
	}
	return mod.DuoGuiItems{
		Slug:  folder,
		Items: items,
	}
}

func (d *DuoVUEdb) DbReadAllComponents() map[string]mod.DuoVUEcomp {
	componentsRaw, err := d.DB.ReadAll("components")
	if err != nil {
		fmt.Println("Error", err)
	}
	components := make(map[string]mod.DuoVUEcomp)
	for _, componentRaw := range componentsRaw {
		component := mod.DuoVUEcomp{}
		if err := json.Unmarshal([]byte(componentRaw), &component); err != nil {
			fmt.Println("Error", err)
		}
		components[component.ID] = component
	}
	return components
}

func (d *DuoVUEdb) DbReadAddressBook() map[string]string {
	addressbookRaw, err := d.DB.ReadAll("addressbook")
	if err != nil {
		fmt.Println("Error", err)
	}
	addressbook := make(map[string]string)
	for _, addressRaw := range addressbookRaw {
		address := mod.AddBook{}
		if err := json.Unmarshal([]byte(addressRaw), &address); err != nil {
			fmt.Println("Error", err)
		}
		addressbook[address.Address] = address.Label
	}
	return addressbook
}
func (d *DuoVUEdb) DbRead(folder, name string) {
	item := mod.DuoGuiItem{}
	if err := d.DB.Read(folder, name, &item); err != nil {
		fmt.Println("Error", err)
	}
	d.Data = item
	fmt.Println("Daasdddddddaaa", item)
}
func (d *DuoVUEdb) DbWrite(folder, name string, data interface{}) {
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
