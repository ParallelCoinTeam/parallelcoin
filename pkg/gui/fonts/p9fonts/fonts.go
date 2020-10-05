package p9fonts

import (
	"fmt"
	"sync"

	"gioui.org/font/opentype"
	"gioui.org/text"

	"github.com/p9c/pod/pkg/gui/fonts/bariolbold"
	"github.com/p9c/pod/pkg/gui/fonts/bariolbolditalic"
	"github.com/p9c/pod/pkg/gui/fonts/bariollight"
	"github.com/p9c/pod/pkg/gui/fonts/bariollightitalic"
	"github.com/p9c/pod/pkg/gui/fonts/bariolregular"
	"github.com/p9c/pod/pkg/gui/fonts/bariolregularitalic"
	"github.com/p9c/pod/pkg/gui/fonts/plan9"
)

var (
	once       sync.Once
	collection []text.FontFace
)

func Collection() []text.FontFace {
	once.Do(func() {
		register(text.Font{Typeface: "bariol"}, bariolregular.TTF)
		register(text.Font{Typeface: "plan9"}, plan9.TTF)
		register(text.Font{Typeface: "bariol", Style: text.Italic}, bariolregularitalic.TTF)
		register(text.Font{Typeface: "bariol", Weight: text.Bold}, bariolbold.TTF)
		register(text.Font{Typeface: "bariol", Style: text.Italic, Weight: text.Bold}, bariolbolditalic.TTF)
		register(text.Font{Typeface: "bariol", Weight: text.Medium}, bariollight.TTF)
		register(text.Font{Typeface: "bariol", Weight: text.Medium, Style: text.Italic}, bariollightitalic.TTF)
		// register(text.Font{Typeface: "go"}, gomono.TTF)
		// register(text.Font{Typeface: "go", Weight: text.Bold}, gomonobold.TTF)
		// register(text.Font{Typeface: "go", Weight: text.Bold, Style: text.Italic}, gomonobolditalic.TTF)
		// register(text.Font{Typeface: "go", Style: text.Italic}, gomonoitalic.TTF)
		// register(text.Font{Typeface: "go", Style: text.Italic}, gomonoitalic.TTF)
		// Ensure that any outside appends will not reuse the backing store.
		n := len(collection)
		collection = collection[:n:n]
	})
	return collection
}

func register(fnt text.Font, ttf []byte) {
	face, err := opentype.Parse(ttf)
	if err != nil {
		panic(fmt.Errorf("failed to parse font: %v", err))
	}
	fnt.Typeface = "Go"
	collection = append(collection, text.FontFace{Font: fnt, Face: face})
}
