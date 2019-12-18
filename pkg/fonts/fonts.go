package fonts

import (
	"fmt"
	"github.com/p9c/gio-parallel/font"
	"github.com/p9c/gio-parallel/font/opentype"
	"github.com/p9c/gio-parallel/text"
	"github.com/p9c/pod/pkg/fonts/bariolbold"
	"github.com/p9c/pod/pkg/fonts/bariolbolditalic"
	"github.com/p9c/pod/pkg/fonts/bariollightitalic"
	"github.com/p9c/pod/pkg/fonts/bariollight"
	"github.com/p9c/pod/pkg/fonts/bariolregular"
	"github.com/p9c/pod/pkg/fonts/bariolregularitalic"
)

func Register() {
	register(text.Font{}, bariolregular.TTF)
	register(text.Font{Style: text.Italic}, bariolregularitalic.TTF)
	register(text.Font{Weight: text.Bold}, bariolbold.TTF)
	register(text.Font{Style: text.Italic, Weight: text.Bold}, bariolbolditalic.TTF)
	register(text.Font{Weight: text.Medium}, bariollight.TTF)
	register(text.Font{Weight: text.Medium, Style: text.Italic}, bariollightitalic.TTF)
}

func register(fnt text.Font, ttf []byte) {
	face, err := opentype.Parse(ttf)
	if err != nil {
		panic(fmt.Sprintf("failed to parse font: %v", err))
	}
	fnt.Typeface = "Go"
	font.Register(fnt, face)
}
