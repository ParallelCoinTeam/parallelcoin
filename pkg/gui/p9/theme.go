package p9

import (
	"gioui.org/text"
	"gioui.org/unit"
	qu "github.com/p9c/pod/pkg/util/quit"
)

type Theme struct {
	quit          qu.C
	shaper        text.Shaper
	collection    []text.FontFace
	TextSize      unit.Value
	Colors        Colors
	icons         map[string]*Icon
	scrollBarSize int
	Dark          *bool
	iconCache     IconCache
	WidgetPool    *Pool
}

// NewTheme creates a new theme to use for rendering a user interface
func NewTheme(fontCollection []text.FontFace, quit qu.C) (th *Theme) {
	th = &Theme{
		quit:          quit,
		shaper:        text.NewCache(fontCollection),
		collection:    fontCollection,
		TextSize:      unit.Sp(16),
		Colors:        NewColors(),
		scrollBarSize: 0,
		iconCache:     make(IconCache),
	}
	th.WidgetPool = th.NewPool()
	return
}
