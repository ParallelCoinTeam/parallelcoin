package gui

import (
	"regexp"
	
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"
	
	l "gioui.org/layout"
	
	"github.com/atotto/clipboard"
)

type Input struct {
	*Window
	editor               *Editor
	input                *TextInput
	clearClickable       *Clickable
	clearButton          *IconButton
	copyClickable        *Clickable
	copyButton           *IconButton
	pasteClickable       *Clickable
	pasteButton          *IconButton
	GetText              func() string
	SetText              func(string)
	SetPasteFunc         func() bool
	borderColor          string
	borderColorUnfocused string
	borderColorFocused   string
	backgroundColor      string
	focused              bool
}

var findSpaceRegexp = regexp.MustCompile(`\s+`)

func (w *Window) Input(txt, hint, borderColorFocused, borderColorUnfocused,
	backgroundColor string, handle func(txt string), ) *Input {
	editor := w.Editor().SingleLine().Submit(true)
	input := w.TextInput(editor, hint).TextScale(1)
	p := &Input{
		Window:               w,
		clearClickable:       w.Clickable(),
		copyClickable:        w.Clickable(),
		pasteClickable:       w.Clickable(),
		editor:               editor,
		input:                input,
		borderColorUnfocused: borderColorUnfocused,
		borderColorFocused:   borderColorFocused,
		backgroundColor:      backgroundColor,
	}
	p.GetText = func() string {
		return p.editor.Text()
	}
	p.SetText = func(s string) {
		p.editor.SetText(s)
	}
	p.clearButton = w.IconButton(p.clearClickable)
	p.copyButton = w.IconButton(p.copyClickable)
	p.pasteButton = w.IconButton(p.pasteClickable)
	clearClickableFn := func() {
		p.editor.SetText("")
		p.editor.Focus()
	}
	copyClickableFn := func() {
		if err := clipboard.WriteAll(p.editor.Text()); Check(err) {
		}
		p.editor.Focus()
	}
	pasteClickableFn := func() {
		col := p.editor.Caret.Col
		if p.SetPasteFunc != nil {
			if p.SetPasteFunc() {
				return
			}
		}
		txt := p.editor.Text()
		var err error
		var cb string
		if cb, err = clipboard.ReadAll(); Check(err) {
		}
		// the SetPasteFunc is a function that screens the clipboard input and
		// does something other than directly paste into the field. If it returns
		// true, this means it has captured input that is processed other than as
		// a simple paste
		cb = findSpaceRegexp.ReplaceAllString(cb, " ")
		txt = txt[:col] + cb + txt[col:]
		p.editor.SetText(txt)
		p.editor.Move(col + len(cb))
		p.editor.Focus()
	}
	p.clearButton.
		Icon(
			w.Icon().
				Color("DocText").
				Src(&icons2.ContentBackspace),
		)
	p.copyButton.
		Icon(
			w.Icon().
				Color("DocText").
				Src(&icons2.ContentContentCopy),
		)
	p.pasteButton.
		Icon(
			w.Icon().
				Color("DocText").
				Src(&icons2.ContentContentPaste),
		)
	p.input.Color("DocText")
	p.clearClickable.SetClick(clearClickableFn)
	p.copyClickable.SetClick(copyClickableFn)
	p.pasteClickable.SetClick(pasteClickableFn)
	p.editor.SetText(txt).SetSubmit(
		func(txt string) {
			go func() {
				handle(txt)
			}()
		},
	).SetChange(
		func(txt string) {
			// send keystrokes to the NSA
		},
	)
	p.editor.SetFocus(
		func(is bool) {
			if is {
				p.borderColor = p.borderColorFocused
			} else {
				p.borderColor = p.borderColorUnfocused
			}
		},
	)
	return p
}

// Fn renders the input widget
func (in *Input) Fn(gtx l.Context) l.Dimensions {
	// gtx.Constraints.Max.X = int(in.TextSize.Scale(float32(in.size)).V)
	// gtx.Constraints.Min.X = 0
	// width := int(in.Theme.TextSize.Scale(in.size).V)
	// gtx.Constraints.Max.X, gtx.Constraints.Min.X = width, width
	return in.Border().Width(0.125).CornerRadius(0.0).
		Corners(0).Color(in.borderColor).Embed(
		in.Fill(in.backgroundColor, l.Center, in.TextSize.V, 0,
			in.Inset(0.25,
				in.Flex().
					Flexed(1,
						in.Inset(0.125, in.input.Color("DocText").Fn).Fn,
					).
					Rigid(
						in.copyButton.
							Background("").
							Icon(in.Icon().Color(in.borderColor).Scale(Scales["H6"]).Src(&icons2.ContentContentCopy)).
							ButtonInset(0.25).
							Fn,
					).
					Rigid(
						in.pasteButton.
							Background("").
							Icon(in.Icon().Color(in.borderColor).Scale(Scales["H6"]).Src(&icons2.ContentContentPaste)).
							ButtonInset(0.25).
							Fn,
					).
					Rigid(
						in.clearButton.
							Background("").
							Icon(in.Icon().Color(in.borderColor).Scale(Scales["H6"]).Src(&icons2.ContentBackspace)).
							ButtonInset(0.25).
							Fn,
					).
					Fn,
			).Fn,
		).Fn).Fn(gtx)
}
