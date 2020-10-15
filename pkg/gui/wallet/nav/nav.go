package nav

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/p9c/pod/pkg/gui/wallet/btn"
	"github.com/p9c/pod/pkg/gui/wallet/dap/page"
	"github.com/p9c/pod/pkg/gui/wallet/lyt"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
)

var (
	navList = &layout.List{}
)

type Navigation struct {
	Name         string
	Bg           string
	Logo         Logo
	Wide         bool
	Mode         string
	NavLayout    string
	ItemLayout   string
	ItemIconSize unit.Value
	Axis         layout.Axis
	Size         int
	NoContent    bool
	LogoWidget   layout.Widget
	MenuItems    []Item
	CurrentPage  page.Page
}
type Item struct {
	Title     string
	Bg        string
	Icon      *widget.Icon
	Btn       *widget.Clickable
	Page      page.Page
	HideOnMob bool
}

func (n *Navigation) Nav(th *theme.Theme, gtx C) layout.Widget {
	return func(gtx C) D {
		n.Size = 64
		n.NoContent = true
		if n.Wide {
			n.Size = 64
			n.NoContent = false
		}
		switch n.Mode {
		case "mobile":
			gtx.Constraints.Min.Y = n.Size
			gtx.Constraints.Max.Y = n.Size
		case "tablet":
			gtx.Constraints.Min.X = n.Size
			gtx.Constraints.Max.X = n.Size
		default:
			gtx.Constraints.Min.X = n.Size * 2
			gtx.Constraints.Max.X = n.Size * 2
		}

		navList.Axis = n.Axis
		return lyt.Format(gtx, n.NavLayout,
			n.LogoWidget,
			//func(gtx C)D{ return D{}},
			func(gtx C) D {
				return navList.Layout(gtx, len(n.MenuItems), func(gtx C, i int) D {
					b := D{}
					item := n.MenuItems[i]
					if !item.HideOnMob && !n.Wide {
						btn := btn.IconTextBtn(th, item.Btn, n.ItemLayout, n.NoContent, item.Title)
						btn.TextSize = unit.Dp(12)
						btn.Icon = item.Icon
						btn.IconSize = n.ItemIconSize
						btn.IconColor = th.Colors["NavItem"]
						btn.CornerRadius = unit.Dp(0)
						btn.Background = th.Colors["NavBg"]
						btn.TextColor = th.Colors["NavItem"]
						for item.Btn.Clicked() {
							n.CurrentPage = item.Page
						}
						b = btn.Layout(gtx)
					}
					return b
					//func(gtx C)D{ return D{}})
				})
			})
	}
}
