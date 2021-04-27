package cfg

import (
	"sort"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/gio/text"

	l "github.com/p9c/gio/layout"

	"github.com/p9c/gel"
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
		return ng.Theme.VFlex().Rigid(
			ng.H6(it.label).Fn,
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
	// schema := podcfg.GetConfigSchema(c.cx.Config)
	tabNames := make(GroupsMap)
	// // tabs := make(p9.WidgetMap)
	// for i := range schema.Groups {
	// 	for j := range schema.Groups[i].Fields {
	// 		sgf := schema.Groups[i].Fields[j]
	// 		if _, ok := tabNames[sgf.Group]; !ok {
	// 			tabNames[sgf.Group] = make(ItemMap)
	// 		}
	// 		tabNames[sgf.Group][sgf.Slug] = &Item{
	// 			slug:        sgf.Slug,
	// 			typ:         sgf.Type,
	// 			label:       sgf.Label,
	// 			description: sgf.Title,
	// 			widget:      sgf.Widget,
	// 			dataType:    sgf.Datatype,
	// 			options:     sgf.Options,
	// 			Slot:        c.cx.ConfigMap[sgf.Slug],
	// 		}
	// 		// D.S(sgf)
	// 		// create all the necessary widgets required before display
	// 		tgs := tabNames[sgf.Group][sgf.Slug]
	// 		switch sgf.Widget {
	// 		case "toggle":
	// 			c.Bools[sgf.Slug] = c.Bool(*tgs.Slot.(*bool)).SetOnChange(
	// 				func(b bool) {
	// 					D.Ln(sgf.Slug, "submitted", b)
	// 					bb := c.cx.ConfigMap[sgf.Slug].(*bool)
	// 					*bb = b
	// 					podcfg.Save(c.cx.Config)
	// 					if sgf.Slug == "DarkTheme" {
	// 						c.Theme.Colors.SetTheme(b)
	// 					}
	// 				},
	// 			)
	// 		case "integer":
	// 			c.inputs[sgf.Slug] = c.Input(
	// 				fmt.Sprint(*tgs.Slot.(*int)), sgf.Slug, "DocText", "DocBg", "PanelBg", func(txt string) {
	// 					D.Ln(sgf.Slug, "submitted", txt)
	// 					i := c.cx.ConfigMap[sgf.Slug].(*int)
	// 					if n, e := strconv.Atoi(txt); !E.Chk(e) {
	// 						*i = n
	// 					}
	// 					podcfg.Save(c.cx.Config)
	// 				}, nil,
	// 			)
	// 		case "time":
	// 			c.inputs[sgf.Slug] = c.Input(
	// 				fmt.Sprint(
	// 					*tgs.Slot.(*time.
	// 					Duration),
	// 				), sgf.Slug, "DocText", "DocBg", "PanelBg", func(txt string) {
	// 					D.Ln(sgf.Slug, "submitted", txt)
	// 					tt := c.cx.ConfigMap[sgf.Slug].(*time.Duration)
	// 					if d, e := time.ParseDuration(txt); !E.Chk(e) {
	// 						*tt = d
	// 					}
	// 					podcfg.Save(c.cx.Config)
	// 				}, nil,
	// 			)
	// 		case "float":
	// 			c.inputs[sgf.Slug] = c.Input(
	// 				strconv.FormatFloat(
	// 					*tgs.Slot.(
	// 					*float64), 'f', -1, 64,
	// 				), sgf.Slug, "DocText", "DocBg", "PanelBg", func(txt string) {
	// 					D.Ln(sgf.Slug, "submitted", txt)
	// 					ff := c.cx.ConfigMap[sgf.Slug].(*float64)
	// 					if f, e := strconv.ParseFloat(txt, 64); !E.Chk(e) {
	// 						*ff = f
	// 					}
	// 					podcfg.Save(c.cx.Config)
	// 				}, nil,
	// 			)
	// 		case "string":
	// 			c.inputs[sgf.Slug] = c.Input(
	// 				*tgs.Slot.(*string), sgf.Slug, "DocText", "DocBg", "PanelBg", func(txt string) {
	// 					D.Ln(sgf.Slug, "submitted", txt)
	// 					ss := c.cx.ConfigMap[sgf.Slug].(*string)
	// 					*ss = txt
	// 					podcfg.Save(c.cx.Config)
	// 				}, nil,
	// 			)
	// 		case "password":
	// 			c.passwords[sgf.Slug] = c.Password(
	// 				"password",
	// 				tgs.Slot.(*string), "DocText", "DocBg", "PanelBg",
	// 				func(txt string) {
	// 					D.Ln(sgf.Slug, "submitted", txt)
	// 					pp := c.cx.ConfigMap[sgf.Slug].(*string)
	// 					*pp = txt
	// 					podcfg.Save(c.cx.Config)
	// 				},
	// 			)
	// 		case "multi":
	// 			c.multis[sgf.Slug] = c.Multiline(
	// 				tgs.Slot.(*cli.StringSlice), "DocText", "DocBg", "PanelBg", 30, func(txt []string) {
	// 					D.Ln(sgf.Slug, "submitted", txt)
	// 					sss := c.cx.ConfigMap[sgf.Slug].(*cli.StringSlice)
	// 					*sss = txt
	// 					podcfg.Save(c.cx.Config)
	// 				},
	// 			)
	// 			// c.multis[sgf.Slug]
	// 		case "radio":
	// 			c.checkables[sgf.Slug] = c.Checkable()
	// 			for i := range sgf.Options {
	// 				c.checkables[sgf.Slug+sgf.Options[i]] = c.Checkable()
	// 			}
	// 			txt := *tabNames[sgf.Group][sgf.Slug].Slot.(*string)
	// 			c.enums[sgf.Slug] = c.Enum().SetValue(txt).SetOnChange(
	// 				func(value string) {
	// 					rr := c.cx.ConfigMap[sgf.Slug].(*string)
	// 					*rr = value
	// 					podcfg.Save(c.cx.Config)
	// 				},
	// 			)
	// 			c.lists[sgf.Slug] = c.List()
	// 		}
	// 	}
	// }

	// D.S(tabNames)
	return tabNames // .Widget(c)
	// return func(gtx l.Context) l.Dimensions {
	// 	return l.Dimensions{}
	// }
}

func (gm GroupsMap) Widget(ng *Config) l.Widget {
	// _, file, line, _ := runtime.Caller(2)
	// D.F("%s:%d", file, line)
	var groups Lists
	for i := range gm {
		var li ListItems
		gmi := gm[i]
		for j := range gmi {
			gmij := gmi[j]
			li = append(
				li, ListItem{
					name: j,
					widget: func() []l.Widget {
						return ng.RenderConfigItem(gmij, len(li))
					},
					// },
				},
			)
		}
		sort.Sort(li)
		groups = append(groups, List{name: i, items: li})
	}
	sort.Sort(groups)
	var out []l.Widget
	first := true
	for i := range groups {
		// D.Ln(groups[i].name)
		g := groups[i]
		if !first {
			// put a space between the sections
			out = append(
				out, func(gtx l.Context) l.Dimensions {
					dims := ng.VFlex().
						// Rigid(
						// 	// ng.Inset(0.25,
						// 	ng.Fill("DocBg", l.Center, ng.TextSize.True, l.S, ng.Inset(0.25,
						// 		gel.EmptyMaxWidth()).Fn,
						// 	).Fn,
						// 	// ).Fn,
						// ).
						Rigid(ng.Inset(0.25, gel.EmptyMaxWidth()).Fn).
						// Rigid(
						// 	// ng.Inset(0.25,
						// 	ng.Fill("DocBg", l.Center, ng.TextSize.True, l.N, ng.Inset(0.25,
						// 		gel.EmptyMaxWidth()).Fn,
						// 	).Fn,
						// 	// ).Fn,
						// ).
						Fn(gtx)
					// ng.Fill("PanelBg", gel.EmptySpace(gtx.Constraints.Max.X, gtx.Constraints.Max.Y), l.Center, 0).Fn(gtx)
					return dims
				},
			)
			// out = append(out, func(gtx l.Context) l.Dimensions {
			// 	return ng.th.ButtonInset(0.25, p9.EmptySpace(0, 0)).Fn(gtx)
			// })
		} else {
			first = false
		}
		// put in the header
		out = append(
			out,
			ng.Fill(
				"Primary", l.Center, ng.TextSize.V*2, 0, ng.Flex().Flexed(
					1,
					ng.Inset(
						0.75,
						ng.H3(g.name).
							Color("DocText").
							Alignment(text.Start).
							Fn,
					).Fn,
				).Fn,
			).Fn,
		)
		// out = append(out, func(gtx l.Context) l.Dimensions {
		// 	return ng.th.Fill("PanelBg",
		// 		ng.th.ButtonInset(0.25,
		// 			ng.th.Flex().Flexed(1,
		// 				p9.EmptyMaxWidth(),
		// 			).Fn,
		// 		).Fn,
		// 	).Fn(gtx)
		// })
		// add the widgets
		for j := range groups[i].items {
			gi := groups[i].items[j]
			for x := range gi.widget() {
				k := x
				out = append(
					out, func(gtx l.Context) l.Dimensions {
						if k < len(gi.widget()) {
							return ng.Fill(
								"DocBg", l.Center, ng.TextSize.V, 0, ng.Flex().
									// Rigid(
									// 	ng.Inset(0.25, gel.EmptySpace(0, 0)).Fn,
									// ).
									Rigid(
										ng.Inset(
											0.25,
											gi.widget()[k],
										).Fn,
									).Fn,
							).Fn(gtx)
						}
						return l.Dimensions{}
					},
				)
			}
		}
	}
	le := func(gtx l.Context, index int) l.Dimensions {
		return out[index](gtx)
	}
	return func(gtx l.Context) l.Dimensions {
		// clip.UniformRRect(f32.Rectangle{
		// 	Max: f32.Pt(float32(gtx.Constraints.Max.X), float32(gtx.Constraints.Max.Y)),
		// }, ng.TextSize.True/2).Add(gtx.Ops)
		return ng.Fill(
			"DocBg", l.Center, ng.TextSize.V, 0, ng.Inset(
				0.25,
				ng.lists["settings"].
					Vertical().
					Length(len(out)).
					// Background("PanelBg").
					// Color("DocBg").
					// Active("Primary").
					ListElement(le).
					Fn,
			).Fn,
		).Fn(gtx)
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
	D.Ln("fallthrough", item.widget)
	return []l.Widget{func(l.Context) l.Dimensions { return l.Dimensions{} }}
}

func (c *Config) RenderToggle(item *Item) []l.Widget {
	return []l.Widget{
		func(gtx l.Context) l.Dimensions {
			return c.Inset(
				0.25,
				c.Flex().
					Rigid(
						c.Switch(c.Bools[item.slug]).DisabledColor("Light").Fn,
					).
					Flexed(
						1,
						c.VFlex().
							Rigid(
								c.Body1(item.label).Fn,
							).
							Rigid(
								c.Caption(item.description).Fn,
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
			return c.Inset(
				0.25,
				c.Flex().Flexed(
					1,
					c.VFlex().
						Rigid(
							c.Body1(item.label).Fn,
						).
						Rigid(
							c.inputs[item.slug].Fn,
						).
						Rigid(
							c.Caption(item.description).Fn,
						).
						Fn,
				).Fn,
			).Fn(gtx)
		},
	}
}

func (c *Config) RenderTime(item *Item) []l.Widget {
	return []l.Widget{
		func(gtx l.Context) l.Dimensions {
			return c.Inset(
				0.25,
				c.Flex().Flexed(
					1,
					c.VFlex().
						Rigid(
							c.Body1(item.label).Fn,
						).
						Rigid(
							c.inputs[item.slug].Fn,
						).
						Rigid(
							c.Caption(item.description).Fn,
						).
						Fn,
				).Fn,
			).
				Fn(gtx)
		},
	}
}

func (c *Config) RenderFloat(item *Item) []l.Widget {
	return []l.Widget{
		func(gtx l.Context) l.Dimensions {
			return c.Inset(
				0.25,
				c.Flex().Flexed(
					1,
					c.VFlex().
						Rigid(
							c.Body1(item.label).Fn,
						).
						Rigid(
							c.inputs[item.slug].Fn,
						).
						Rigid(
							c.Caption(item.description).Fn,
						).
						Fn,
				).Fn,
			).
				Fn(gtx)
		},
	}
}

func (c *Config) RenderString(item *Item) []l.Widget {
	return []l.Widget{
		c.Inset(
			0.25,
			c.Flex().Flexed(
				1,
				c.VFlex().
					Rigid(
						c.Body1(item.label).Fn,
					).
					Rigid(
						c.inputs[item.slug].Fn,
					).
					Rigid(
						c.Caption(item.description).Fn,
					).
					Fn,
			).Fn,
		).
			Fn,
	}
}

func (c *Config) RenderPassword(item *Item) []l.Widget {
	return []l.Widget{
		c.Inset(
			0.25,
			c.Flex().Flexed(
				1,
				c.VFlex().
					Rigid(
						c.Body1(item.label).Fn,
					).
					Rigid(
						c.passwords[item.slug].Fn,
					).
					Rigid(
						c.Caption(item.description).Fn,
					).
					Fn,
			).Fn,
		).
			Fn,
	}
}

func (c *Config) RenderMulti(item *Item, position int) []l.Widget {
	// D.Ln("rendering multi")
	// c.multis[item.slug].
	w := []l.Widget{
		func(gtx l.Context) l.Dimensions {
			return c.Inset(
				0.25,
				c.Flex().Flexed(
					1,
					c.VFlex().
						Rigid(
							c.Body1(item.label).Fn,
						).
						Rigid(
							c.Caption(item.description).Fn,
						).Fn,
				).Fn,
			).
				Fn(gtx)
		},
	}
	widgets := c.multis[item.slug].Widgets()
	// D.Ln(widgets)
	w = append(w, widgets...)
	return w
}

func (c *Config) RenderRadio(item *Item) []l.Widget {
	out := func(gtx l.Context) l.Dimensions {
		var options []l.Widget
		for i := range item.options {
			var color string
			color = "PanelBg"
			if c.enums[item.slug].Value() == item.options[i] {
				color = "Primary"
			}
			options = append(
				options,
				c.RadioButton(
					c.checkables[item.slug+item.options[i]].
						Color("DocText").
						IconColor(color).
						CheckedStateIcon(&icons.ToggleRadioButtonChecked).
						UncheckedStateIcon(&icons.ToggleRadioButtonUnchecked),
					c.enums[item.slug], item.options[i], item.options[i],
				).Fn,
			)
		}
		return c.Inset(
			0.25,
			c.VFlex().
				Rigid(
					c.Body1(item.label).Fn,
				).
				Rigid(
					c.Flex().
						Rigid(
							func(gtx l.Context) l.Dimensions {
								gtx.Constraints.Max.X = int(c.Theme.TextSize.Scale(10).V)
								return c.lists[item.slug].DisableScroll(true).Slice(gtx, options...)(gtx)
								// 	// return c.lists[item.slug].Length(len(options)).Vertical().ListElement(func(gtx l.Context, index int) l.Dimensions {
								// 	// 	return options[index](gtx)
								// 	// }).Fn(gtx)
								// 	return c.lists[item.slug].Slice(gtx, options...)(gtx)
								// 	// return l.Dimensions{}
							},
						).
						Flexed(
							1,
							c.Caption(item.description).Fn,
						).
						Fn,
				).Fn,
		).
			Fn(gtx)
	}
	return []l.Widget{out}
}
