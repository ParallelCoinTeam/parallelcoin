package fonts

import (
	"fmt"
	"github.com/p9c/pod/pkg/fonts/bariolbold"
	"github.com/p9c/pod/pkg/fonts/bariolbolditalic"
	"github.com/p9c/pod/pkg/fonts/bariollight"
	"github.com/p9c/pod/pkg/fonts/bariollightitalic"
	"github.com/p9c/pod/pkg/fonts/bariolregular"
	"github.com/p9c/pod/pkg/fonts/bariolregularitalic"
	"github.com/p9c/pod/pkg/fonts/plan9"
	"github.com/p9c/pod/pkg/gui/font"
	"github.com/p9c/pod/pkg/gui/font/opentype"
	"github.com/p9c/pod/pkg/gui/text"
)


func Register() {
	register(text.Font{Typeface: "plan9"}, plan9.TTF)
	register(text.Font{Typeface: "bariol"}, bariolregular.TTF)
	register(text.Font{Typeface: "bariol",Style: text.Italic}, bariolregularitalic.TTF)
	register(text.Font{Typeface: "bariol",Weight: text.Bold}, bariolbold.TTF)
	register(text.Font{Typeface: "bariol",Style: text.Italic, Weight: text.Bold}, bariolbolditalic.TTF)
	register(text.Font{Typeface: "bariol",Weight: text.Medium}, bariollight.TTF)
	register(text.Font{Typeface: "bariol",Weight: text.Medium, Style: text.Italic}, bariollightitalic.TTF)
}

func register(fnt text.Font, ttf []byte) {
	face, err := opentype.Parse(ttf)
	if err != nil {
		panic(fmt.Sprintf("failed to parse font: %v", err))
	}
	font.Register(fnt, face)
}
