package cfg

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	l "gioui.org/layout"
	"github.com/urfave/cli"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/app/save"
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
	Slot        interface{}
}

func (it *Item) Item(ng *Config) l.Widget {
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

func (c *Config) Config() GroupsMap {
	schema := pod.GetConfigSchema(c.cx.Config, c.cx.ConfigMap)
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
				Slot:        c.cx.ConfigMap[sgf.Slug],
			}
			// Debugs(sgf)
			// create all the necessary widgets required before display
			tgs := tabNames[sgf.Group][sgf.Slug]
			switch sgf.Widget {
			case "toggle":
				c.Bools[sgf.Slug] = c.th.Bool(*tgs.Slot.(*bool)).SetOnChange(func(b bool) {
					Debug(sgf.Slug, "submitted", b)
					bb := c.cx.ConfigMap[sgf.Slug].(*bool)
					*bb = b
					save.Pod(c.cx.Config)
					if sgf.Slug == "DarkTheme" {
						c.th.Colors.SetTheme(b)
					}
				})
			case "integer":
				c.inputs[sgf.Slug] = c.th.Input(fmt.Sprint(*tgs.Slot.(*int)), sgf.Slug,
					"Primary", "PanelBg", 26, func(txt string) {
						Debug(sgf.Slug, "submitted", txt)
						i := c.cx.ConfigMap[sgf.Slug].(*int)
						if n, err := strconv.Atoi(txt); !Check(err) {
							*i = n
						}
						save.Pod(c.cx.Config)
					})
			case "time":
				c.inputs[sgf.Slug] = c.th.Input(fmt.Sprint(*tgs.Slot.(*time.Duration)), sgf.Slug,
					"Primary", "PanelBg", 26, func(txt string) {
						Debug(sgf.Slug, "submitted", txt)
						tt := c.cx.ConfigMap[sgf.Slug].(*time.Duration)
						if d, err := time.ParseDuration(txt); !Check(err) {
							*tt = d
						}
						save.Pod(c.cx.Config)
					})
			case "float":
				c.inputs[sgf.Slug] = c.th.Input(strconv.FormatFloat(*tgs.Slot.(*float64), 'f', -1, 64), sgf.Slug,
					"Primary", "PanelBg", 26, func(txt string) {
						Debug(sgf.Slug, "submitted", txt)
						ff := c.cx.ConfigMap[sgf.Slug].(*float64)
						if f, err := strconv.ParseFloat(txt, 64); !Check(err) {
							*ff = f
						}
						save.Pod(c.cx.Config)
					})
			case "string":
				c.inputs[sgf.Slug] = c.th.Input(*tgs.Slot.(*string), sgf.Slug,
					"Primary", "PanelBg", 26, func(txt string) {
						Debug(sgf.Slug, "submitted", txt)
						ss := c.cx.ConfigMap[sgf.Slug].(*string)
						*ss = txt
						save.Pod(c.cx.Config)
					})
			case "password":
				c.passwords[sgf.Slug] = c.th.Password("password", tgs.Slot.(*string),
					"Primary", "PanelBg", 26, func(txt string) {
						Debug(sgf.Slug, "submitted", txt)
						pp := c.cx.ConfigMap[sgf.Slug].(*string)
						*pp = txt
						save.Pod(c.cx.Config)
					})
			case "multi":
				c.multis[sgf.Slug] = c.th.Multiline(tgs.Slot.(*cli.StringSlice),
					"Primary", "PanelBg", 30, func(txt []string) {
						Debug(sgf.Slug, "submitted", txt)
						sss := c.cx.ConfigMap[sgf.Slug].(*cli.StringSlice)
						*sss = txt
						save.Pod(c.cx.Config)
					})
				// c.multis[sgf.Slug]
			case "radio":
				c.checkables[sgf.Slug] = c.th.Checkable()
				for i := range sgf.Options {
					c.checkables[sgf.Slug+sgf.Options[i]] = c.th.Checkable()
				}
				txt := *tabNames[sgf.Group][sgf.Slug].Slot.(*string)
				c.enums[sgf.Slug] = c.th.Enum().SetValue(txt).SetOnChange(func(value string) {
					rr := c.cx.ConfigMap[sgf.Slug].(*string)
					*rr = value
					save.Pod(c.cx.Config)
				})
				c.lists[sgf.Slug] = c.th.List()
			}
		}
	}

	// Debugs(tabNames)
	return tabNames // .Widget(c)
	// return func(gtx l.Context) l.Dimensions {
	// 	return l.Dimensions{}
	// }
}

func (gm GroupsMap) Widget(ng *Config) l.Widget {
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
				widget: func() []l.Widget {
					return ng.RenderConfigItem(gmij, len(li))
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
			// out = append(out, func(gtx l.Context) l.Dimensions {
			// 	return ng.th.Fill("DocBg",
			// 		ng.th.Inset(0.25,
			// 			ng.th.Flex().Flexed(1, p9.EmptySpace(0, 0)).Fn,
			// 		).Fn,
			// 	).Fn(gtx)
			// })
			out = append(out, func(gtx l.Context) l.Dimensions {
				return ng.th.Inset(0.25, p9.EmptySpace(0, 0)).Fn(gtx)
			})
		} else {
			first = false
		}
		// put in the header
		out = append(out, func(gtx l.Context) l.Dimensions {
			return ng.th.Inset(0.0, ng.th.H6(g.name).Color("PanelText").Fn).Fn(gtx)
		})
		out = append(out, func(gtx l.Context) l.Dimensions {
			return ng.th.Fill("DocBg",
				ng.th.Inset(0.25,
					ng.th.Flex().Flexed(1, p9.EmptySpace(0, 0)).Fn,
				).Fn,
			).Fn(gtx)
		})
		// add the widgets
		for j := range groups[i].items {
			gi := groups[i].items[j]
			for x := range gi.widget() {
				k := x
				out = append(out, func(gtx l.Context) l.Dimensions {
					// return ng.th.Fill("DocBg",
					// 	ng.th.Inset(0.25,
					// 		return func(gtx l.Context) l.Dimensions {
					if k < len(gi.widget()) {

						return ng.th.Fill("DocBg",
							ng.th.Flex().
								Rigid(
									ng.th.Inset(0.25, p9.EmptySpace(0, 0)).Fn,
								).
								Rigid(
									gi.widget()[k],
								).Fn,
						).Fn(gtx)
					}
					return l.Dimensions{}
					// }
					// ).Fn,
					// ).Fn(gtx)
				})
			}
		}
	}
	le := func(gtx l.Context, index int) l.Dimensions {
		return out[index](gtx)
	}
	return func(gtx l.Context) l.Dimensions {
		return ng.th.Inset(0.25, ng.lists["settings"].
			Vertical().
			Length(len(out)).
			Background("PanelBg").
			// Color("DocText").Active("Primary").
			ListElement(le).Fn).Fn(gtx)
	}
}

// RenderConfigItem renders a config item. It takes a position variable which tells it which index it begins on
// the bigger config widget list, with this and its current data set the multi can insert and delete elements above
// its add button without rerendering the config item or worse, the whole config widget
func (c *Config) RenderConfigItem(item *Item, position int) []l.Widget {
	switch item.widget {
	case "toggle":
		return c.RenderToggle(item)
	case "integer":
		return c.RenderInteger(item)
	case "time":
		return c.RenderTime(item)
	case "float":
		return c.RenderFloat(item)
	case "string":
		return c.RenderString(item)
	case "password":
		return c.RenderPassword(item)
	case "multi":
		return c.RenderMulti(item, position)
	case "radio":
		return c.RenderRadio(item)
	}
	Debug("fallthrough", item.widget)
	return []l.Widget{func(l.Context) l.Dimensions { return l.Dimensions{} }}
}

func (c *Config) RenderToggle(item *Item) []l.Widget {
	return []l.Widget{
		func(gtx l.Context) l.Dimensions {
			return c.th.Inset(0.25, c.th.Flex().
				Rigid(
					c.th.Switch(c.Bools[item.slug]).Fn,
				).
				Rigid(
					c.th.VFlex().
						Rigid(
							c.th.Body1(item.label).Fn,
						).
						Rigid(
							c.th.Caption(item.description).Fn,
						).
						Fn,
				).Fn,
			).Fn(gtx)
		},
	}
}

func (c *Config) RenderInteger(item *Item) []l.Widget {
	return []l.Widget{
		func(gtx l.Context) l.Dimensions {
			return c.th.Inset(0.25, c.th.VFlex().
				Rigid(
					c.th.Body1(item.label).Fn,
				).
				Rigid(
					c.inputs[item.slug].Fn,
				).
				Rigid(
					c.th.Caption(item.description).Fn,
				).
				Fn,
			).
				Fn(gtx)
		},
	}
}

func (c *Config) RenderTime(item *Item) []l.Widget {
	return []l.Widget{
		func(gtx l.Context) l.Dimensions {
			return c.th.Inset(0.25, c.th.VFlex().
				Rigid(
					c.th.Body1(item.label).Fn,
				).
				Rigid(
					c.inputs[item.slug].Fn,
				).
				Rigid(
					c.th.Caption(item.description).Fn,
				).
				Fn,
			).
				Fn(gtx)
		},
	}
}

func (c *Config) RenderFloat(item *Item) []l.Widget {
	return []l.Widget{
		func(gtx l.Context) l.Dimensions {
			return c.th.Inset(0.25, c.th.VFlex().
				Rigid(
					c.th.Body1(item.label).Fn,
				).
				Rigid(
					c.inputs[item.slug].Fn,
				).
				Rigid(
					c.th.Caption(item.description).Fn,
				).
				Fn,
			).
				Fn(gtx)
		},
	}
}

func (c *Config) RenderString(item *Item) []l.Widget {
	return []l.Widget{
		c.th.Inset(0.25,
			c.th.VFlex().
				Rigid(
					c.th.Body1(item.label).Fn,
				).
				Rigid(
					c.inputs[item.slug].Fn,
				).
				Rigid(
					c.th.Caption(item.description).Fn,
				).
				Fn,
		).
			Fn,
	}
}

func (c *Config) RenderPassword(item *Item) []l.Widget {
	return []l.Widget{
		c.th.Inset(0.25,
			c.th.VFlex().
				Rigid(
					c.th.Body1(item.label).Fn,
				).
				Rigid(
					c.passwords[item.slug].Fn,
				).
				Rigid(
					c.th.Caption(item.description).Fn,
				).
				Fn,
		).
			Fn,
	}
}

func (c *Config) RenderMulti(item *Item, position int) []l.Widget {
	// Debug("rendering multi")
	// c.multis[item.slug].
	w := []l.Widget{
		func(gtx l.Context) l.Dimensions {
			return c.th.Inset(0.25,
				c.th.VFlex().
					Rigid(
						c.th.Body1(item.label).Fn,
					).
					Rigid(
						c.th.Caption(item.description).Fn,
					).Fn,
			).
				Fn(gtx)
		},
	}
	widgets := c.multis[item.slug].Widgets()
	// Debug(widgets)
	w = append(w, widgets...)
	return w
}

func (c *Config) RenderRadio(item *Item) []l.Widget {
	out := func(gtx l.Context) l.Dimensions {
		var options []l.Widget
		for i := range item.options {
			color := "DocText"
			if c.enums[item.slug].Value() == item.options[i] {
				color = "Primary"
			}
			options = append(options,
				c.th.RadioButton(
					c.checkables[item.slug+item.options[i]].
						IconColor(color).
						Color(color).
						CheckedStateIcon(&icons.ToggleRadioButtonChecked).
						UncheckedStateIcon(&icons.ToggleRadioButtonUnchecked),
					c.enums[item.slug], item.options[i], item.options[i]).Fn)
		}
		return c.th.Inset(0.25,
			c.th.VFlex().
				Rigid(
					c.th.Body1(item.label).Fn,
				).
				Rigid(
					c.th.Flex().
						Rigid(
							func(gtx l.Context) l.Dimensions {
								gtx.Constraints.Max.X = int(c.th.TextSize.Scale(10).V)
								return c.lists[item.slug].DisableScroll(true).Slice(gtx, options...)(gtx)
								// 	// return c.lists[item.slug].Length(len(options)).Vertical().ListElement(func(gtx l.Context, index int) l.Dimensions {
								// 	// 	return options[index](gtx)
								// 	// }).Fn(gtx)
								// 	return c.lists[item.slug].Slice(gtx, options...)(gtx)
								// 	// return l.Dimensions{}
							},
						).
						Rigid(
							c.th.Caption(item.description).Fn,
						).
						Fn,
				).Fn,
		).
			Fn(gtx)
	}
	return []l.Widget{out}
}
