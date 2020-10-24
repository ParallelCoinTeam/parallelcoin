package p9

import (
	l "gioui.org/layout"
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"
)

type Input struct {
	*Theme
	editor               *Editor
	input                *TextInput
	clearClickable       *Clickable
	clearButton          *IconButton
	GetText              func() string
	size                 int
	borderColor          string
	borderColorUnfocused string
	borderColorFocused   string
	focused              bool
}

func (th *Theme) Input(txt, borderColorFocused, borderColorUnfocused string, size int, handle func(txt string)) *Input {
	editor := th.Editor().SingleLine().Submit(true)
	input := th.SimpleInput(editor)
	p := &Input{
		Theme:          th,
		clearButton:    nil,
		clearClickable: th.Clickable(),
		editor:         editor,
		input:          input,
		size:           size,
		borderColorUnfocused: borderColorUnfocused,
		borderColorFocused: borderColorFocused,
	}
	p.GetText = func() string {
		return p.editor.Text()
	}
	p.clearButton = th.IconButton(p.clearClickable)
	clearClickableFn := func() {
		p.editor.SetText("")
	}
	p.clearButton.
		Icon(
			th.Icon().
				Color("DocText").
				Src(icons2.ContentBackspace),
		)
	p.input.Color("DocText")
	p.clearClickable.SetClick(clearClickableFn)
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
				in.clearButton.
					Background("").
					Icon(in.Icon().Color(in.borderColor).Src(icons2.ContentBackspace)).Fn,
			).
			Fn,
	).Fn(gtx)
}
