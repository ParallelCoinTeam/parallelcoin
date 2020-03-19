package monitor

import (
	"os"
	"sort"
	"strings"
)

type Node struct {
	Name     string
	Parent   *Node
	Children []*Node
}

func GetTree(paths []string) (root *Node) {
	sort.Strings(paths)
	root = &Node{}
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
				Name:     name,
				Parent:   cursor,
				Children: nil,
			}
			cursor.Children = append(cursor.Children, n)
		case splitLen == prevLen:
			// attach sibling
			n = &Node{
				Name:     name,
				Parent:   cursor.Parent,
				Children: nil,
			}
			cursor.Parent.Children = append(cursor.Parent.Children, n)
		case splitLen < prevLen:
			cursor = cursor.Parent
			n = &Node{
				Name:     name,
				Parent:   cursor.Parent,
				Children: nil,
			}
			cursor.Parent.Children = append(cursor.Parent.Children, n)
		}
		cursor = n
		prevLen = splitLen
	}
	// spew.Config.Indent = "    "
	// L.Debugs(root)
	return
}
