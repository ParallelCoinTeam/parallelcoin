package monitor

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/util/logi/Pkg/Pk"
	"github.com/p9c/pod/pkg/util/logi/consume"
	"os"
	"sort"
	"strings"
)

type Node struct {
	Name               string
	FullName           string
	parent             *Node
	Children           []*Node
	Closed             bool
	Hidden             bool
	foldButton         *gel.Button
	showButton         *gel.Button
	showChildrenButton *gel.Button
	hideChildrenButton *gel.Button
	empty              bool
}

func (s *State) GetTree(paths []string) (root *Node) {
	sort.Strings(paths)
	var sliced [][]string
	for i := range paths {
		sliced = append(sliced, strings.Split(paths[i], "/"))
	}
	slicedPaths := make(map[string]bool)
	for i := range sliced {
		var s string
		for j := range sliced[i] {
			empty := true
			if j == len(sliced[i])+1 {
				empty = false
			} else {
				s = strings.Join(sliced[i][:j+1], "/")
				//Debug(s)
			}
			slicedPaths[s] = empty
		}
	}
	paths = make([]string, len(slicedPaths))
	counter := 0
	for i := range slicedPaths {
		paths[counter] = i
		counter++
	}
	Debugs(slicedPaths)
	sort.Strings(paths)
	Debugs(paths)
	s.FilterRoot = &Node{
		Name:       "root",
		FullName:   string(os.PathSeparator),
		parent:     nil,
		Children:   nil,
		Closed:     false,
		Hidden:     false,
		showButton: nil,
	}
	cursor := s.FilterRoot
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
				Name:               name,
				FullName:           v,
				parent:             cursor,
				foldButton:         new(gel.Button),
				showButton:         new(gel.Button),
				showChildrenButton: new(gel.Button),
				hideChildrenButton: new(gel.Button),
				Closed:             false,
				Hidden:             false,
			}
			cursor.Children = append(cursor.Children, n)
		case splitLen == prevLen:
			// attach sibling
			n = &Node{
				Name:               name,
				FullName:           v,
				parent:             cursor.parent,
				foldButton:         new(gel.Button),
				showButton:         new(gel.Button),
				showChildrenButton: new(gel.Button),
				hideChildrenButton: new(gel.Button),
				Closed:             false,
				Hidden:             false,
			}
			cursor.parent.Children = append(cursor.parent.Children, n)
		case splitLen < prevLen:
			cursor = cursor.parent
			n = &Node{
				Name:               name,
				FullName:           v,
				parent:             cursor.parent,
				foldButton:         new(gel.Button),
				showButton:         new(gel.Button),
				showChildrenButton: new(gel.Button),
				hideChildrenButton: new(gel.Button),
				Closed:             false,
				Hidden:             false,
			}
			cursor.parent.Children = append(cursor.parent.Children, n)
		default:
			n = &Node{
				Closed: false,
				Hidden: false,
			}
		}
		s.Config.FilterNodes[v] = n
		cursor = n
		prevLen = splitLen
	}
	//root.ClearParents()
	//root.CloseAllItems(s)
	//spew.Config.Indent = "    "
	//Debugs(root)
	return s.FilterRoot
}

func (n *Node) GetWidget(s *State, headless bool) {
	gtx := s.Gtx
	if headless {
		gtx = s.Htx
	}
	nn := n.GetOpenItems()[1:]
	indent := 0
	s.Lists["Filter"].Axis = layout.Vertical
	s.Lists["Filter"].Layout(gtx, len(nn), func(i int) {
		s.FlexH(
			gui.Rigid(func() {
				split := strings.Split(nn[i].FullName, string(os.PathSeparator))
				indent = len(split) - 1
				s.Inset(0, func() {
					s.Rectangle(indent*16, 32, "PanelBg", "ff")
				})
			}),
			gui.Rigid(func() {
				name := nn[i].Name
				if name == "" {
					name = "root"
				}
				fg := "PanelText"
				if nn[i].Hidden {
					fg = "DocBg"
				}
				s.TextButton(name, "Primary", 24, fg, "PanelBg", nn[i].showButton)
				if nn[i].showButton.Clicked(gtx) {
					nn[i].Hidden = !nn[i].Hidden
					//if !nn[i].Hidden {
					//	nn[i].ShowAllItems(s)
					//} else {
					//	nn[i].HideAllItems(s)
					//}
					consume.SetFilter(s.Worker, s.FilterRoot.GetPackages())
					s.SaveConfig()
				}
			}),
			gui.Rigid(func() {
				if len(nn[i].Children) > 0 {
					ic := "Folded"
					if !nn[i].Closed {
						ic = "Unfolded"
					}
					fg := "DocBg"
					if nn[i].IsAnyShowing() {
						fg = "DocTextDim"
					}
					if nn[i].IsAllShowing() {
						fg = "PanelText"
					}
					s.IconButton(ic, fg, "PanelBg", nn[i].foldButton)
					if nn[i].foldButton.Clicked(gtx) {
						if nn[i].Closed {
							//nn[i].OpenAllItems(s)
						} else {
							nn[i].CloseAllItems(s)
						}
						nn[i].Closed = !nn[i].Closed
						s.SaveConfig()
					}
				}
			}),
			gui.Flexed(1, func() {

			}),
			gui.Rigid(func() {
				if len(nn[i].Children) > 0 && nn[i].IsAnyHiding() {
					s.IconButton("ShowItem", "DocBg", "PanelBg",
						nn[i].showChildrenButton)
					for nn[i].showChildrenButton.Clicked(gtx) {
						Debug("filter all")
						nn[i].ShowAllItems(s)
						nn[i].Hidden = false
						consume.SetFilter(s.Worker, s.FilterRoot.GetPackages())
						s.SaveConfig()
					}
				}
			}), gui.Rigid(func() {
				if len(nn[i].Children) > 0 && nn[i].IsAnyShowing() {
					s.IconButton("HideItem", "DocBg", "PanelBg",
						nn[i].hideChildrenButton)
					for nn[i].hideChildrenButton.Clicked(gtx) {
						Debug("filter none")
						nn[i].Hidden = true
						nn[i].HideAllItems(s)
						consume.SetFilter(s.Worker, s.FilterRoot.GetPackages())
						s.SaveConfig()
					}
				}
			}),
		)
	})
}

func (n *Node) GetPackages() (out Pk.Package) {
	out = make(Pk.Package)
	all := n.GetAllItems()
	for i := range all {
		out[all[i].FullName] = !all[i].Hidden
	}
	return
}

func (n *Node) GetAllItems() (out []*Node) {
	out = append(out, n)
	for _, v := range n.Children {
		out = append(out, v.GetAllItems()...)
	}
	return
}

func (n *Node) GetOpenItems() (out []*Node) {
	out = append(out, n)
	for _, v := range n.Children {
		if !n.Closed {
			out = append(out, v.GetOpenItems()...)
		}
	}
	return
}

func (n *Node) OpenAllItems(s *State) {
	for _, v := range n.Children {
		v.Closed = false
		v.OpenAllItems(s)
	}
}

func (n *Node) CloseAllItems(s *State) {
	for _, v := range n.Children {
		v.Closed = true
		v.CloseAllItems(s)
	}
}

func (n *Node) HideAllItems(s *State) {
	for _, v := range n.Children {
		v.Hidden = true
		v.HideAllItems(s)
	}
}

func (n *Node) ShowAllItems(s *State) {
	for _, v := range n.Children {
		v.Hidden = false
		v.ShowAllItems(s)
	}
}

func (n *Node) IsAnyShowing() bool {
	for _, v := range n.Children {
		if !v.Hidden || v.IsAnyShowing() {
			return true
		}
	}
	return false
}

func (n *Node) IsAnyHiding() bool {
	for _, v := range n.Children {
		if v.Hidden || v.IsAnyHiding() {
			return true
		}
	}
	return false
}

func (n *Node) IsAllShowing() bool {
	for _, v := range n.Children {
		if v.Hidden || !v.IsAllShowing() {
			return false
		}
	}
	return true
}

func (n *Node) IsNoneShowing() bool {
	for _, v := range n.Children {
		if !v.Hidden || v.IsNoneShowing() {
			return false
		}
	}
	return true
}
