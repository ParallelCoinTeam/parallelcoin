package podcfg

import (
	"github.com/p9c/pod/pkg/logg"
	"reflect"
	"sort"
)

type Schema struct {
	Groups Groups `json:"groups"`
}
type Groups []Group

type Group struct {
	Legend string `json:"legend"`
	Fields `json:"fields"`
}

type Fields []Field

func (f Fields) Len() int {
	return len(f)
}

func (f Fields) Less(i, j int) bool {
	return f[i].Label < f[j].Label
}

func (f Fields) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

type Field struct {
	Datatype    string   `json:"datatype"`
	Description string   `json:"help"`
	Featured    string   `json:"featured"`
	Group       string   `json:"group"`
	Hooks       string   `json:"hooks"`
	Label       string   `json:"label"`
	Model       string   `json:"model"`
	Options     []string `json:"options"`
	Restart     string   `json:"restart"`
	Slug        string   `json:"slug"`
	Type        string   `json:"type"`
	Widget      string   `json:"inputType"`
}

// GetConfigSchema returns a schema for a given config
func GetConfigSchema(cfg *Config) Schema {
	t := reflect.TypeOf(cfg)
	t = t.Elem()
	var levelOptions, network []string
	for i := range logg.LevelSpecs {
		levelOptions = append(levelOptions, logg.LevelSpecs[i].Name)
	}
	network = []string{"mainnet", "testnet", "regtestnet", "simnet"}
	rawFields := make(map[string]Fields)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		var options []string
		switch {
		case field.Name == "LogLevel":
			options = levelOptions
		case field.Name == "Network":
			options = network
		}
		f := Field{
			Datatype:    field.Type.String(),
			Description: field.Tag.Get("description"),
			Featured:    field.Tag.Get("featured"),
			Group:       field.Tag.Get("group"),
			Hooks:       field.Tag.Get("hooks"),
			Label:       field.Tag.Get("label"),
			Model:       field.Tag.Get("json"),
			Options:     options,
			Slug:        field.Name,
			Type:        field.Tag.Get("type"),
			Widget:      field.Tag.Get("widget"),
		}
		if f.Group != "" {
			rawFields[f.Group] = append(rawFields[f.Group], f)
		}
	}
	for i := range rawFields {
		sort.Sort(rawFields[i])
	}
	var outGroups Groups
	var rf []string
	for i := range rawFields {
		rf = append(rf, i)
	}
	sort.Strings(rf)
	for i := range rf {
		rf[i], rf[len(rf)-1-i] = rf[len(rf)-1-i], rf[i]
	}
	for i := range rf {
		group := Group{
			Legend: rf[i],
			Fields: rawFields[rf[i]],
		}
		outGroups = append(Groups{group}, outGroups...)
	}
	return Schema{
		Groups: outGroups,
	}
}

