package monitor

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gel"
	"os"
	"sort"
	"strings"
)

type Node struct {
	Name     string
	FullName string
	parent   *Node
	Children []*Node
	Closed   bool
	Hidden   bool
	//Node       *TreeNode
	foldButton *gel.Button
	showButton *gel.Button
}

func (s *State) GetTree(paths []string) (root *Node) {

	sort.Strings(paths)
	root = &Node{
		Name:       "root",
		FullName:   string(os.PathSeparator),
		parent:     nil,
		Children:   nil,
		Closed:     false,
		Hidden:     false,
		foldButton: new(gel.Button),
		showButton: nil,
	}
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
				parent:     cursor,
				foldButton: new(gel.Button),
				showButton: new(gel.Button),
				Closed:     false,
				Hidden:     false,
			}
			cursor.Children = append(cursor.Children, n)
		case splitLen == prevLen:
			// attach sibling
			n = &Node{
				Name:       name,
				FullName:   v,
				parent:     cursor.parent,
				foldButton: new(gel.Button),
				showButton: new(gel.Button),
				Closed:     false,
				Hidden:     false,
			}
			cursor.parent.Children = append(cursor.parent.Children, n)
		case splitLen < prevLen:
			cursor = cursor.parent
			n = &Node{
				Name:       name,
				FullName:   v,
				parent:     cursor.parent,
				foldButton: new(gel.Button),
				showButton: new(gel.Button),
				Closed:     false,
				Hidden:     false,
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
	//L.Debugs(root)
	return
}

func (n *Node) GetWidget(s *State) {
	nn := n.GetOpenItems()[1:]
	//for i := range nn {
	//	L.Debug(nn[i].FullName, nn[i].Closed, nn[i].Hidden)
	//}
	s.FilterList.Axis = layout.Vertical
	s.FilterList.Layout(s.Gtx, len(nn), func(i int) {
		s.FlexH(
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
				if nn[i].Hidden {
					fg = "DocBg"
				}
				s.TextButton(name, "Primary", 24, fg, "PanelBg", nn[i].showButton)
				if nn[i].showButton.Clicked(s.Gtx) {
					nn[i].Hidden = !nn[i].Hidden
					if !nn[i].Hidden {
						nn[i].ShowAllItems(s)
					} else {
						nn[i].HideAllItems(s)
					}

					s.SaveConfig()
				}
			}),
			Rigid(func() {
				if len(nn[i].Children) > 0 {
					ic := "Folded"
					if !nn[i].Closed {
						ic = "Unfolded"
					}
					fg := "PanelText"
					if !nn[i].IsAnyShowing() {
						fg = "DocBg"
					}
					s.IconButton(ic, fg, "PanelBg", nn[i].foldButton)
					if nn[i].foldButton.Clicked(s.Gtx) {
						nn[i].Closed = !nn[i].Closed
					}

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

//func (n *Node) LoadState(s *State) {
//	//s.Config.FilterNodes = make(map[string]TreeNode)
//	for _, v := range n.Children {
//		s.Config.FilterNodes[v.FullName] = &TreeNode{
//			Closed: n.Closed,
//			Hidden: n.Hidden,
//		}
//		v.LoadState(s)
//	}
//}
//
//func (n *Node) StoreState(s *State) {
//	for j, w := range n.Children {
//		s.Config.FilterNodes[n.Children[j].FullName].Closed =
//			w.Closed
//		s.Config.FilterNodes[n.Children[j].FullName].Hidden =
//			w.Hidden
//		w.StoreState(s)
//	}
//}

//func (n *Node) ClearParents() {
//	for i := range n.Children {
//		n.Children[i].parent = nil
//		n.Children[i].ClearParents()
//	}
//}

func (n *Node) GetOpenItems() (out []*Node) {
	out = append(out, n)
	//L.Debugs(n)
	for _, v := range n.Children {
		if !n.Closed {
			out = append(out, v.GetOpenItems()...)
		}
	}
	//if n.Parent == nil {
	//	L.Debugs(out)
	//}
	return
}

func (n *Node) OpenAllItems(s *State) {
	for _, v := range n.Children {
		v.Closed = false
		//s.Config.FilterNodes[v.FullName].Node.Closed=false
		v.OpenAllItems(s)
	}
}

func (n *Node) CloseAllItems(s *State) {
	for _, v := range n.Children {
		v.Closed = true
		//s.Config.FilterNodes[v.FullName].Closed=true
		v.CloseAllItems(s)
	}
}

func (n *Node) HideAllItems(s *State) {
	for _, v := range n.Children {
		v.Hidden = true
		//s.Config.FilterNodes[v.FullName].Hidden=true
		v.HideAllItems(s)
	}
}

func (n *Node) ShowAllItems(s *State) {
	for _, v := range n.Children {
		v.Hidden = false
		//s.Config.FilterNodes[v.FullName].Hidden=false
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
