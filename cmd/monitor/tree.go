package monitor

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gel"
	"os"
	"sort"
	"strings"
)

type Node struct {
	Name       string
	FullName   string
	Parent     *Node
	Children   []*Node
	Opened     bool
	FoldButton *gel.Button
	Hide       bool
	ShowButton *gel.Button
}

func GetTree(paths []string) (root *Node) {
	sort.Strings(paths)
	root = &Node{FoldButton: new(gel.Button)}
	cursor := root
	var prevLen int
	for _, v := range paths {
		split := strings.Split(v, string(os.PathSeparator))
		splitLen := len(split)
		name := split[splitLen-1]
		var n *Node
		switch {
		case splitLen > prevLen:
			// attach child
			n = &Node{
				Name:       name,
				FullName:   v,
				Parent:     cursor,
				FoldButton: new(gel.Button),
				ShowButton: new(gel.Button),
			}
			cursor.Children = append(cursor.Children, n)
		case splitLen == prevLen:
			// attach sibling
			n = &Node{
				Name:       name,
				FullName:   v,
				Parent:     cursor.Parent,
				FoldButton: new(gel.Button),
				ShowButton: new(gel.Button),
			}
			cursor.Parent.Children = append(cursor.Parent.Children, n)
		case splitLen < prevLen:
			cursor = cursor.Parent
			n = &Node{
				Name:       name,
				FullName:   v,
				Parent:     cursor.Parent,
				FoldButton: new(gel.Button),
				ShowButton: new(gel.Button),
			}
			cursor.Parent.Children = append(cursor.Parent.Children, n)
		}
		cursor = n
		prevLen = splitLen
	}
	//root.CloseAllItems()
	// spew.Config.Indent = "    "
	// L.Debugs(root)
	return
}

func (n *Node) GetWidget(s *State) {
	//L.Debug("drawing filter list")
	s.Loggers.LoadState(s)
	nn := n.GetOpenItems()[1:]
	s.FilterList.Axis = layout.Vertical
	s.FilterList.Layout(s.Gtx, len(nn), func(i int) {
		s.FlexH(
			//Rigid(func() {
			//	ic := "ShowItem"
			//	if nn[i].Show.Load() {
			//		ic = "HideItem"
			//	}
			//	s.IconButton(ic, "DocBg", "PanelBg", nn[i].ShowButton)
			//
			//}),
			Rigid(func() {
				split := strings.Split(nn[i].FullName, string(os.PathSeparator))
				joined := strings.Join(split[:len(split)-1], string(os.PathSeparator))
				if joined != "" {
					s.Text(joined+string(os.PathSeparator), "DocBg", "PanelBg", "Primary", "body2")()
				}
			}),
			Rigid(func() {
				name := nn[i].Name
				if name == "" {
					name = "root"
				}
				fg := "PanelText"
				if !nn[i].Hide {
					fg = "DocBg"
				}
				//s.Text(name, fg, "PanelBg", "Primary", "h6")()
				s.TextButton(name, "Primary", 24, fg, "PanelBg", nn[i].ShowButton)
				if nn[i].ShowButton.Clicked(s.Gtx) {
					nn[i].Hide = !nn[i].Hide
					if !nn[i].Hide {
						nn[i].ShowAllItems()
					} else {
						nn[i].HideAllItems()
					}

					nn[i].StoreState(s)
					s.SaveConfig()
				}
			}),
			Rigid(func() {
				if len(nn[i].Children) > 0 {
					ic := "Unfolded"
					if nn[i].Opened {
						ic = "Folded"
					}
					fg := "PanelText"
					if nn[i].IsAnyShowing() {
						fg = "DocBg"
					}
					s.IconButton(ic, fg, "PanelBg", nn[i].FoldButton)
					if nn[i].FoldButton.Clicked(s.Gtx) {
						nn[i].Opened = !nn[i].Opened
					}
					nn[i].StoreState(s)
					s.SaveConfig()
				}
			}),
			Flexed(1, func() {

			}),
			//Rigid(func() {
			//	s.Label("x")
			//}),
		)
	})
}

func (n *Node) LoadState(s *State) {
	//s.Config.FilterTreeNodes = make(map[string]TreeNode)
	for _, v := range n.Children {
		s.Config.FilterTreeNodes[v.FullName] = &TreeNode{
			IsOpen: n.Opened,
			Hidden: n.Hide,
		}
		v.LoadState(s)
	}
}

func (n *Node) StoreState(s *State) {
	for j, w := range n.Children {
		s.Config.FilterTreeNodes[n.Children[j].FullName].IsOpen =
			w.Opened
		s.Config.FilterTreeNodes[n.Children[j].FullName].Hidden =
			w.Hide
		w.StoreState(s)
	}
}

func (n *Node) GetOpenItems() (out []*Node) {
	out = append(out, n)
	for _, v := range n.Children {
		if !n.Opened {
			out = append(out, v.GetOpenItems()...)
		}
	}
	return
}

func (n *Node) OpenAllItems() {
	for _, v := range n.Children {
		v.Opened = true
		v.OpenAllItems()
	}
}

func (n *Node) CloseAllItems() {
	for _, v := range n.Children {
		v.Opened = false
		v.CloseAllItems()
	}
}

func (n *Node) HideAllItems() {
	for _, v := range n.Children {
		v.Hide = true
		v.HideAllItems()
	}
}

func (n *Node) ShowAllItems() {
	for _, v := range n.Children {
		v.Hide = false
		v.ShowAllItems()
	}
}

func (n *Node) IsAnyShowing() bool {
	for _, v := range n.Children {
		if !v.Hide || !v.IsAnyShowing() {
			return true
		}
	}
	return false
}
