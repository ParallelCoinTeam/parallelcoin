package gui

import (
	"gioui.org/text"
	"gioui.org/unit"
	
	qu "github.com/p9c/pod/pkg/util/quit"
)

type Theme struct {
	quit       qu.C
	shaper     text.Shaper
	collection Collection
	TextSize   unit.Value
	*Colors
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
		Colors:        newColors(),
		scrollBarSize: 0,
		iconCache:     make(IconCache),
	}
	return
}
