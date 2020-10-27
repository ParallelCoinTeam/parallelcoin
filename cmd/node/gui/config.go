package gui

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	l "gioui.org/layout"
	"github.com/urfave/cli"
	"golang.org/x/exp/shiny/materialdesign/icons"

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
		return ng.th.VFlex().Rigid(
			ng.th.H6(it.label).Fn,
		).Fn(gtx)
	}
}

type ItemMap map[string]*Item

type GroupsMap map[string]ItemMap

type ListItem struct {
	name   string
	widget func() []l.Widget
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

func (ng *NodeGUI) Config() GroupsMap {
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
			tgs := tabNames[sgf.Group][sgf.Slug]
			switch sgf.Widget {
			case "toggle":
				ng.bools[sgf.Slug] = ng.th.Bool(*tgs.slot.(*bool))
			case "integer":
				ng.inputs[sgf.Slug] = ng.th.Input(fmt.Sprint(*tgs.slot.(*int)),
					"Primary", "PanelBg", 26, func(txt string) {
						Debug(sgf.Slug, "submitted", txt)
					})
			case "time":
				ng.inputs[sgf.Slug] = ng.th.Input(fmt.Sprint(*tgs.slot.(*time.Duration)),
					"Primary", "PanelBg", 26, func(txt string) {
						Debug(sgf.Slug, "submitted", txt)
					})
			case "float":
				ng.inputs[sgf.Slug] = ng.th.Input(strconv.FormatFloat(*tgs.slot.(*float64), 'f', -1, 64),
					"Primary", "PanelBg", 26, func(txt string) {
						Debug(sgf.Slug, "submitted", txt)
					})
			case "string":
				ng.inputs[sgf.Slug] = ng.th.Input(*tgs.slot.(*string),
					"Primary", "PanelBg", 26, func(txt string) {
						Debug(sgf.Slug, "submitted", txt)
					})
			case "password":
				ng.passwords[sgf.Slug] = ng.th.Password(tgs.slot.(*string),
					"Primary", "PanelBg", 26, func(txt string) {
						Debug(sgf.Slug, "submitted", txt)
					})
			case "multi":
				ng.multis[sgf.Slug] = ng.th.Multiline(tgs.slot.(*cli.StringSlice),
					"Primary", "PanelBg", 30, func(txt []string) {
						Debug(sgf.Slug, "submitted", txt)
					})
				// ng.multis[sgf.Slug]
			case "radio":
				ng.checkables[sgf.Slug] = ng.th.Checkable()
				for i := range sgf.Options {
					ng.checkables[sgf.Slug+sgf.Options[i]] = ng.th.Checkable()
				}
				ng.enums[sgf.Slug] = ng.th.Enum().SetValue(*tabNames[sgf.Group][sgf.Slug].slot.(*string))
				ng.lists[sgf.Slug] = ng.th.List()
			}
		}
	}

	// Debugs(tabNames)
	return tabNames // .Widget(ng)
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
				widget:
				func() []l.Widget {
					return ng.RenderConfigItem(gmij)
				},
				// },
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
				return ng.th.Inset(0.5, p9.EmptySpace(0, 0)).Fn(gtx)
			})
		} else {
			first = false
		}
		// put in the header
		out = append(out, func(gtx l.Context) l.Dimensions {
			return ng.th.Inset(0.0, ng.th.Fill("DocText", ng.th.Inset(0.5, ng.th.H6(g.name).Color("DocBg").Fn).Fn).Fn).Fn(gtx)
		})
		// add the widgets
		for j := range groups[i].items {
			gi := groups[i].items[j]
			for x := range gi.widget() {
				k := x
				out = append(out, func(gtx l.Context) l.Dimensions {
					return ng.th.Fill("DocBg",
						ng.th.Inset(0.25,
							func(gtx l.Context) l.Dimensions {
								return gi.widget()[k](gtx)
							},
						).Fn,
					).Fn(gtx)
				})
			}
		}
	}
	le := func(gtx l.Context, index int) l.Dimensions {
		return out[index](gtx)
	}
	return func(gtx l.Context) l.Dimensions {
		return ng.lists["settings"].Vertical().Length(len(out)).ListElement(le).Fn(gtx)
	}
}

func (ng *NodeGUI) RenderConfigItem(item *Item) []l.Widget {
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
	return []l.Widget{func(l.Context) l.Dimensions { return l.Dimensions{} }}
}

func (ng *NodeGUI) RenderToggle(item *Item) []l.Widget {
	return []l.Widget{
		func(gtx l.Context) l.Dimensions {
			return ng.th.Flex().
				Rigid(
					ng.th.Switch(ng.bools[item.slug]).Fn,
					// p9.EmptySpace(0, 0),
				).
				Rigid(
					ng.th.VFlex().
						Rigid(
							ng.th.Body1(item.label).Fn,
						).
						Rigid(
							ng.th.Caption(item.description).Fn,
						).Fn,
				).Fn(gtx)
		},
	}
}

func (ng *NodeGUI) RenderInteger(item *Item) []l.Widget {
	return []l.Widget{
		func(gtx l.Context) l.Dimensions {
			return ng.th.VFlex().
				Rigid(
					ng.th.Body1(item.label).Fn,
				).
				Rigid(
					ng.inputs[item.slug].Fn,
				).
				Rigid(
					ng.th.Caption(item.description).Fn,
				).
				Fn(gtx)
		},
	}
}

func (ng *NodeGUI) RenderTime(item *Item) []l.Widget {
	return []l.Widget{
		func(gtx l.Context) l.Dimensions {
			return ng.th.VFlex().
				Rigid(
					ng.th.Body1(item.label).Fn,
				).
				Rigid(
					ng.inputs[item.slug].Fn,
				).
				Rigid(
					ng.th.Caption(item.description).Fn,
				).
				Fn(gtx)
		},
	}
}

func (ng *NodeGUI) RenderFloat(item *Item) []l.Widget {
	return []l.Widget{
		func(gtx l.Context) l.Dimensions {
			return ng.th.VFlex().
				Rigid(
					ng.th.Body1(item.label).Fn,
				).
				Rigid(
					ng.inputs[item.slug].Fn,
				).
				Rigid(
					ng.th.Caption(item.description).Fn,
				).
				Fn(gtx)
		},
	}
}

func (ng *NodeGUI) RenderString(item *Item) []l.Widget {
	return []l.Widget{
		func(gtx l.Context) l.Dimensions {
			return ng.th.VFlex().
				Rigid(
					ng.th.Body1(item.label).Fn,
				).
				Rigid(
					ng.inputs[item.slug].Fn,
				).
				Rigid(
					ng.th.Caption(item.description).Fn,
				).
				Fn(gtx)
		},
	}
}

func (ng *NodeGUI) RenderPassword(item *Item) []l.Widget {
	return []l.Widget{
		func(gtx l.Context) l.Dimensions {
			return ng.th.VFlex().
				Rigid(
					ng.th.Body1(item.label).Fn,
				).
				Rigid(
					ng.passwords[item.slug].Fn,
				).
				Rigid(
					ng.th.Caption(item.description).Fn,
				).
				Fn(gtx)
		},
	}
}

func (ng *NodeGUI) RenderMulti(item *Item) []l.Widget {
	// Debug("rendering multi")
	w := []l.Widget{
		func(gtx l.Context) l.Dimensions {
			return ng.th.VFlex().
				Rigid(
					ng.th.Body1(item.label).Fn,
				).
				Rigid(
					ng.th.Caption(item.description).Fn,
				).
				// Rigid(
				// 	ng.multis[item.slug].Fn,
				// ).
				Fn(gtx)
		},
	}
	widgets := ng.multis[item.slug].Widgets()
	// Debug(widgets)
	w = append(w, widgets...)
	return w
}

func (ng *NodeGUI) RenderRadio(item *Item) []l.Widget {
	out := func(gtx l.Context) l.Dimensions {
		var options []l.Widget
		for i := range item.options {
			color := "DocText"
			if ng.enums[item.slug].Value() == item.options[i] {
				color = "Primary"
			}
			options = append(options,
				ng.th.RadioButton(
					ng.checkables[item.slug+item.options[i]].
						IconColor(color).
						Color(color).
						CheckedStateIcon(icons.ToggleRadioButtonChecked).
						UncheckedStateIcon(icons.ToggleRadioButtonUnchecked),
					ng.enums[item.slug], item.options[i], item.options[i]).Fn)
		}
		return ng.th.VFlex().
			Rigid(
				ng.th.Body1(item.label).Fn,
			).
			Rigid(
				ng.th.Flex().
					Rigid(
						func(gtx l.Context) l.Dimensions {
							gtx.Constraints.Max.X = int(ng.th.TextSize.Scale(10).V)
							return ng.lists[item.slug].DisableScroll(true).Slice(gtx, options...)(gtx)
							// 	// return ng.lists[item.slug].Length(len(options)).Vertical().ListElement(func(gtx l.Context, index int) l.Dimensions {
							// 	// 	return options[index](gtx)
							// 	// }).Fn(gtx)
							// 	return ng.lists[item.slug].Slice(gtx, options...)(gtx)
							// 	// return l.Dimensions{}
						},
					).
					Rigid(
						ng.th.Caption(item.description).Fn,
					).
					Fn,
			).
			Fn(gtx)
	}
	return []l.Widget{out}
}
