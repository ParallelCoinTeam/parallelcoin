package gui

import (
	"runtime"
	"sort"

	l "gioui.org/layout"

	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/pod"
)

type Item struct {
	typ         string
	label       string
	description string
	inputType   string
	dataType    string
	options     []string
	slot        interface{}
}

func (it *Item) Item(ng *NodeGUI) l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return ng.VFlex().Rigid(
			ng.H6(it.label).Fn,
		).Fn(gtx)
	}
}

type ItemMap map[string]*Item

type GroupsMap map[string]ItemMap

type ListItem struct {
	name   string
	widget l.Widget
}

type ListItems []ListItem

func (l ListItems) Len() int {
	return len(l)
}

func (l ListItems) Less(i, j int) bool {
	return l[i].name < l[j].name
}

func (l ListItems) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type List struct {
	name  string
	items ListItems
}

type Lists []List

func (l Lists) Len() int {
	return len(l)
}

func (l Lists) Less(i, j int) bool {
	return l[i].name < l[j].name
}

func (l Lists) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (gm GroupsMap) Widget(ng *NodeGUI) l.Widget {
	_, file, line, _ := runtime.Caller(2)
	Debugf("%s:%d", file, line)
	var groups Lists
	for i := range gm {
		var li ListItems
		for j := range gm[i] {
			li = append(li, ListItem{
				name: j,
				widget: func(gtx l.Context) l.Dimensions {
					return ng.H6(gm[i][j].label).Fn(gtx)
				},
			})
		}
		sort.Sort(li)
		groups = append(groups, List{name: i, items: li})
	}
	sort.Sort(groups)
	var out []l.Widget
	first := true
	for i := range groups {
		Debug(groups[i].name)
		g := groups[i]
		if !first {
			out = append(out, func(gtx l.Context) l.Dimensions {
				return ng.Inset(0.5, p9.EmptySpace(0, 0)).Fn(gtx)
			})
		} else {
			first = false
		}
		out = append(out, func(gtx l.Context) l.Dimensions {
			return ng.Inset(0.25, ng.H4(g.name).Fn).Fn(gtx)
		})
		for j := range groups[i].items {
			Debugf("\t%s", groups[i].items[j].name)
			gi := groups[i].items[j]
			out = append(out, func(gtx l.Context) l.Dimensions {
				return ng.Inset(0.5, ng.Body1(gi.name).Fn).Fn(gtx)
			})
		}
	}
	le := func(gtx l.Context, index int) l.Dimensions {
		return out[index](gtx)
		// return l.Dimensions{}
	}
	return func(gtx l.Context) l.Dimensions {
		return ng.lists["settings"].Vertical().Length(len(out)).ListElement(le).Fn(gtx)
	}
}

func (ng *NodeGUI) Config() l.Widget {
	schema := pod.GetConfigSchema(ng.cx.Config, ng.cx.ConfigMap)
	tabNames := make(GroupsMap)
	// tabs := make(p9.WidgetMap)
	for i := range schema.Groups {
		for j := range schema.Groups[i].Fields {
			sgf := schema.Groups[i].Fields[j]
			if _, ok := tabNames[sgf.Group]; !ok {
				tabNames[sgf.Group] = make(ItemMap)
			}
			tabNames[sgf.Group][sgf.Slug] = &Item{
				typ:         sgf.Type,
				label:       sgf.Label,
				description: sgf.Description,
				inputType:   sgf.InputType,
				dataType:    sgf.Datatype,
				options:     sgf.Options,
			}
			// Debugs(sgf)
		}
	}
	// Debugs(tabNames)
	return tabNames.Widget(ng)
	// return func(gtx l.Context) l.Dimensions {
	// 	return l.Dimensions{}
	// }
}
