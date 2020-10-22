package gui

import (
	"sort"

	l "gioui.org/layout"

	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/pod"
)

type Item struct {
	slug        string
	typ         string
	label       string
	description string
	widget      string
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
	// _, file, line, _ := runtime.Caller(2)
	// Debugf("%s:%d", file, line)
	var groups Lists
	for i := range gm {
		var li ListItems
		gmi := gm[i]
		for j := range gmi {
			gmij := gmi[j]
			li = append(li, ListItem{
				name: j,
				widget: func(gtx l.Context) l.Dimensions {
					return ng.RenderConfigItem(gmij)(gtx)
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
		// Debug(groups[i].name)
		g := groups[i]
		if !first {
			// put a space between the sections
			out = append(out, func(gtx l.Context) l.Dimensions {
				return ng.Inset(0.5, p9.EmptySpace(0, 0)).Fn(gtx)
			})
		} else {
			first = false
		}
		// put in the header
		out = append(out, func(gtx l.Context) l.Dimensions {
			return ng.Inset(0.25, ng.H6(g.name).Fn).Fn(gtx)
		})
		// add the widgets
		for j := range groups[i].items {
			// Debugf("\t%s", groups[i].items[j].name)
			gi := groups[i].items[j]
			out = append(out, func(gtx l.Context) l.Dimensions {
				return ng.Inset(0.5,
					gi.widget,
				).Fn(gtx)
			})
		}
	}
	le := func(gtx l.Context, index int) l.Dimensions {
		return out[index](gtx)
	}
	return func(gtx l.Context) l.Dimensions {
		return ng.lists["settings"].Vertical().Length(len(out)).ListElement(le).Fn(gtx)
	}
}

func (ng *NodeGUI) RenderConfigItem(item *Item) l.Widget {
	sl := item.slug
	ty := item.typ
	wi := item.widget
	// dt := item.dataType
	opts := item.options
	slot := item.slot
	Debug(sl, wi, ty, opts, slot)
	return func(gtx l.Context) l.Dimensions {
		return ng.VFlex().
			Rigid(
				ng.Body1(item.label).Fn,
			).
			Rigid(
				ng.Caption(item.description).Fn,
			).
			Fn(gtx)
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
				slug:        sgf.Slug,
				typ:         sgf.Type,
				label:       sgf.Label,
				description: sgf.Description,
				widget:      sgf.Widget,
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
