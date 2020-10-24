package gui

import (
	"fmt"
	"sort"
	"strconv"
	"time"

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
				slot:        ng.cx.ConfigMap[sgf.Slug],
			}
			// Debugs(sgf)
			// create all the necessary widgets required before display
			switch sgf.Widget {
			case "toggle":
				ng.bools[sgf.Slug] = ng.Bool(*tabNames[sgf.Group][sgf.Slug].slot.(*bool))
			case "integer":
				ng.inputs[sgf.Slug] = ng.Input(fmt.Sprint(*tabNames[sgf.Group][sgf.Slug].slot.(*int)),
					"Primary", "PanelBg", 24, func(txt string) {
						Debug(sgf.Slug, "submitted", txt)
					})
			case "time":
				ng.inputs[sgf.Slug] = ng.Input(fmt.Sprint(*tabNames[sgf.Group][sgf.Slug].slot.(*time.Duration)),
					"Primary", "PanelBg", 24, func(txt string) {
						Debug(sgf.Slug, "submitted", txt)
					})
			case "float":
				ng.inputs[sgf.Slug] = ng.Input(strconv.FormatFloat(*tabNames[sgf.Group][sgf.Slug].slot.(*float64), 'f', -1, 64),
					"Primary", "PanelBg", 24, func(txt string) {
						Debug(sgf.Slug, "submitted", txt)
					})
			case "string":
				ng.inputs[sgf.Slug] = ng.Input(*tabNames[sgf.Group][sgf.Slug].slot.(*string),
					"Primary", "PanelBg", 24, func(txt string) {
						Debug(sgf.Slug, "submitted", txt)
					})
			case "password":
				ng.passwords[sgf.Slug] = ng.Password(tabNames[sgf.Group][sgf.Slug].slot.(*string),
					"Primary", "PanelBg", 24, func(txt string) {
						Debug(sgf.Slug, "submitted", txt)
					})
			case "multi":
			case "radio":
				// ng.checkables[sgf.Slug] = ng.Checkable()
				for i := range sgf.Options {
					ng.checkables[sgf.Slug+sgf.Options[i]] = ng.Checkable()
				}
				ng.enums[sgf.Slug] = ng.Enum().SetValue(*tabNames[sgf.Group][sgf.Slug].slot.(*string))
				ng.lists[sgf.Slug] = ng.List()
			}
		}
	}

	// Debugs(tabNames)
	return tabNames.Widget(ng)
	// return func(gtx l.Context) l.Dimensions {
	// 	return l.Dimensions{}
	// }
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
			return ng.Inset(0.0, ng.Fill("DocText", ng.Inset(0.5, ng.H6(g.name).Color("DocBg").Fn).Fn).Fn).Fn(gtx)
		})
		// add the widgets
		for j := range groups[i].items {
			gi := groups[i].items[j]
			out = append(out, func(gtx l.Context) l.Dimensions {
				return ng.Fill("DocBg",
					ng.Inset(0.25,
						gi.widget,
					).Fn,
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
	switch item.widget {
	case "toggle":
		return ng.RenderToggle(item)
	case "integer":
		return ng.RenderInteger(item)
	case "time":
		return ng.RenderTime(item)
	case "float":
		return ng.RenderFloat(item)
	case "string":
		return ng.RenderString(item)
	case "password":
		return ng.RenderPassword(item)
	case "multi":
		return ng.RenderMulti(item)
	case "radio":
		return ng.RenderRadio(item)
	}
	Debug("fallthrough", item.widget)
	return func(l.Context) l.Dimensions { return l.Dimensions{} }
}

func (ng *NodeGUI) RenderToggle(item *Item) l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return ng.Flex().
			Rigid(
				ng.th.Switch(ng.bools[item.slug]).Fn,
				// p9.EmptySpace(0, 0),
			).
			Rigid(
				ng.VFlex().
					Rigid(
						ng.Body1(item.label).Fn,
					).
					Rigid(
						ng.Caption(item.description).Fn,
					).Fn,
			).Fn(gtx)
	}
}

func (ng *NodeGUI) RenderInteger(item *Item) l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return ng.VFlex().
			Rigid(
				ng.Body1(item.label).Fn,
			).
			Rigid(
				ng.inputs[item.slug].Fn,
			).
			Rigid(
				ng.Caption(item.description).Fn,
			).
			Fn(gtx)
	}
}

func (ng *NodeGUI) RenderTime(item *Item) l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return ng.VFlex().
			Rigid(
				ng.Body1(item.label).Fn,
			).
			Rigid(
				ng.inputs[item.slug].Fn,
			).
			Rigid(
				ng.Caption(item.description).Fn,
			).
			Fn(gtx)
	}
}

func (ng *NodeGUI) RenderFloat(item *Item) l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return ng.VFlex().
			Rigid(
				ng.Body1(item.label).Fn,
			).
			Rigid(
				ng.inputs[item.slug].Fn,
			).
			Rigid(
				ng.Caption(item.description).Fn,
			).
			Fn(gtx)
	}
}

func (ng *NodeGUI) RenderString(item *Item) l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return ng.VFlex().
			Rigid(
				ng.Body1(item.label).Fn,
			).
			Rigid(
				ng.inputs[item.slug].Fn,
			).
			Rigid(
				ng.Caption(item.description).Fn,
			).
			Fn(gtx)
	}
}

func (ng *NodeGUI) RenderPassword(item *Item) l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return ng.VFlex().
			Rigid(
				ng.Body1(item.label).Fn,
			).
			Rigid(
				ng.passwords[item.slug].Fn,
			).
			Rigid(
				ng.Caption(item.description).Fn,
			).
			Fn(gtx)
	}
}

func (ng *NodeGUI) RenderMulti(item *Item) l.Widget {
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

func (ng *NodeGUI) RenderRadio(item *Item) l.Widget {
	var options []l.Widget
	for i := range item.options {
		color := "DocText"
		if ng.enums[item.slug].Value() == item.options[i] {
			color = "Primary"
		}
		options = append(options,
			ng.RadioButton(ng.checkables[item.slug+item.options[i]].IconColor(color).Color(color),
				ng.enums[item.slug], item.options[i], item.options[i]).Fn)
	}
	out := func(gtx l.Context) l.Dimensions {
		return ng.VFlex().
			Rigid(
				ng.Body1(item.label).Fn,
			).
			Rigid(
				ng.Flex().
					Rigid(
						func(gtx l.Context) l.Dimensions {
							gtx.Constraints.Max.X = int(ng.TextSize.Scale(10).V)
							// return ng.lists[item.slug].Length(len(options)).Vertical().ListElement(func(gtx l.Context, index int) l.Dimensions {
							// 	return options[index](gtx)
							// }).Fn(gtx)
							return ng.lists[item.slug].Slice(gtx, options...)(gtx)
						},
					).
					Rigid(
						ng.Caption(item.description).Fn,
					).
					Fn,
			).
			Fn(gtx)
	}
	return out
}
