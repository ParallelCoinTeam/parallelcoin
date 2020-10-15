package btn

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gioapp/gel/helper"
	"github.com/p9c/pod/pkg/gui/wallet/lyt"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
	"image"
)

var (
	noReturn = func(gtx C) D { return D{} }
)

type (
	D = layout.Dimensions
	C = layout.Context
)

type IconTextButton struct {
	Theme        *theme.Theme
	Button       *widget.Clickable
	Background   string
	Icon         *widget.Icon
	IconColor    string
	IconSize     unit.Value
	Content      interface{}
	TextSize     unit.Value
	TextColor    string
	CornerRadius unit.Value
	lay          string
	noContent    bool
}

func IconTextBtn(t *theme.Theme, b *widget.Clickable, l string, noContent bool, content interface{}) IconTextButton {
	return IconTextButton{
		Theme:     t,
		Button:    b,
		Icon:      t.Icons["Logo"],
		IconColor: t.Colors["Primary"],
		IconSize:  unit.Dp(16),
		Content:   content,
		TextSize:  unit.Dp(16),
		TextColor: t.Colors["Primary"],
		noContent: noContent,
		lay:       l,
	}
}

func (b IconTextButton) Layout(gtx C) D {
	return material.Clickable(gtx, b.Button, func(gtx C) D {
		var leftMargin float32 = 12
		if b.noContent {
			leftMargin = 8
		}
		return layout.Inset{
			Top:    unit.Dp(8),
			Right:  unit.Dp(8),
			Bottom: unit.Dp(8),
			Left:   unit.Dp(leftMargin),
		}.Layout(gtx, func(gtx C) D {
			layContent := noReturn
			switch c := b.Content.(type) {
			case func(gtx C) D:
				layContent = c
				fmt.Println("test2222", "aAAAA")

			default:
				layContent = b.label(fmt.Sprint(c.(string)))
			}
			if !b.noContent {
				layContent = noReturn
			}
			return lyt.Format(gtx, b.lay,
				b.icon(),
				layContent,
			)
		})
	})
}

func (b IconTextButton) icon() func(gtx C) D {
	return func(gtx C) D {
		return layout.Inset{}.Layout(gtx, func(gtx C) D {
			var d D
			if b.Icon != nil {
				//size := gtx.Px(b.IconSize) - 2*gtx.Px(unit.Dp(16))
				size := gtx.Px(b.IconSize)
				b.Icon.Color = helper.HexARGB(b.IconColor)
				b.Icon.Layout(gtx, unit.Px(float32(size)))
				d = D{
					Size: image.Point{X: size, Y: size},
				}
			}
			return d
		})
	}
}

func (b IconTextButton) label(label string) func(gtx C) D {
	return func(gtx C) D {
		return layout.Inset{
			Top:    unit.Dp(5),
			Right:  unit.Dp(8),
			Bottom: unit.Dp(0),
			Left:   unit.Dp(8),
		}.Layout(gtx, func(gtx C) D {
			l := theme.Body(b.Theme, label)
			l.TextSize = b.TextSize
			l.Alignment = text.Start
			l.Color = helper.HexARGB(b.TextColor)
			return l.Layout(gtx)
		})
	}
}
