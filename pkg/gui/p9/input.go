package p9

import (
	l "gioui.org/layout"
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/clipboard"
)

type Input struct {
	*Theme
	editor               *Editor
	input                *TextInput
	clearClickable       *Clickable
	clearButton          *IconButton
	copyClickable        *Clickable
	copyButton           *IconButton
	pasteClickable       *Clickable
	pasteButton          *IconButton
	GetText              func() string
	size                 int
	borderColor          string
	borderColorUnfocused string
	borderColorFocused   string
	focused              bool
}

func (th *Theme) Input(txt, hint, borderColorFocused, borderColorUnfocused string,
	size int, handle func(txt string)) *Input {
	editor := th.Editor().SingleLine().Submit(true)
	input := th.TextInput(editor, hint)
	p := &Input{
		Theme:                th,
		clearClickable:       th.Clickable(),
		copyClickable:        th.Clickable(),
		pasteClickable:       th.Clickable(),
		editor:               editor,
		input:                input,
		size:                 size,
		borderColorUnfocused: borderColorUnfocused,
		borderColorFocused:   borderColorFocused,
	}
	p.GetText = func() string {
		return p.editor.Text()
	}
	p.clearButton = th.IconButton(p.clearClickable)
	p.copyButton = th.IconButton(p.copyClickable)
	p.pasteButton = th.IconButton(p.pasteClickable)
	clearClickableFn := func() {
		p.editor.SetText("")
		p.editor.Focus()
	}
	copyClickableFn := func() {
		go clipboard.Set(p.editor.Text())
		p.editor.Focus()
	}
	pasteClickableFn := func() {
		go func() {
			txt := p.editor.Text()
			txt = txt[:p.editor.caret.col] + clipboard.Get() + txt[p.editor.caret.col:]
			p.editor.SetText(txt)
		}()
		p.editor.Focus()
	}
	p.clearButton.
		Icon(
			th.Icon().
				Color("DocText").
				Src(&icons2.ContentBackspace),
		)
	p.copyButton.
		Icon(
			th.Icon().
				Color("DocText").
				Src(&icons2.ContentContentCopy),
		)
	p.pasteButton.
		Icon(
			th.Icon().
				Color("DocText").
				Src(&icons2.ContentContentPaste),
		)
	p.input.Color("DocText")
	p.clearClickable.SetClick(clearClickableFn)
	p.copyClickable.SetClick(copyClickableFn)
	p.pasteClickable.SetClick(pasteClickableFn)
	p.editor.SetText(txt).SetSubmit(func(txt string) {
		go func() {
			handle(txt)
		}()
	}).SetChange(func(txt string) {
		// send keystrokes to the NSA
	})
	p.editor.SetFocus(func(is bool) {
		if is {
			p.borderColor = p.borderColorFocused
		} else {
			p.borderColor = p.borderColorUnfocused
		}
	})
	return p
}

func (in *Input) Fn(gtx l.Context) l.Dimensions {
	gtx.Constraints.Max.X = int(in.TextSize.Scale(float32(in.size)).V)
	gtx.Constraints.Min.X = 0
	return in.Border().Color(in.borderColor).Embed(
		in.Flex().
			Flexed(1,
				in.Inset(0.25, in.input.Color("DocText").Fn).Fn,
			).
			Rigid(
				in.copyButton.
					Background("").
					Icon(in.Icon().Color(in.borderColor).Scale(Scales["H6"]).Src(&icons2.ContentContentCopy)).
					Inset(0.25).
					Fn,
			).
			Rigid(
				in.pasteButton.
					Background("").
					Icon(in.Icon().Color(in.borderColor).Scale(Scales["H6"]).Src(&icons2.ContentContentPaste)).
					Inset(0.25).
					Fn,
			).
			Rigid(
				in.clearButton.
					Background("").
					Icon(in.Icon().Color(in.borderColor).Scale(Scales["H6"]).Src(&icons2.ContentBackspace)).
					Inset(0.25).
					Fn,
			).
			Fn,
	).Fn(gtx)
}
