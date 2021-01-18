package gui

import (
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"
	
	l "gioui.org/layout"
	
	"github.com/p9c/pod/pkg/gui/clipboard"
)

type Password struct {
	*Window
	pass                 *Editor
	passInput            *TextInput
	unhideClickable      *Clickable
	unhideButton         *IconButton
	copyClickable        *Clickable
	copyButton           *IconButton
	pasteClickable       *Clickable
	pasteButton          *IconButton
	hide                 bool
	borderColor          string
	borderColorUnfocused string
	borderColorFocused   string
	backgroundColor      string
	focused              bool
	showClickableFn      func(col string)
	password             *string
	handle               func(pass string)
}

func (w *Window) Password(
	hint string, password *string, borderColorFocused,
	borderColorUnfocused, backgroundColor string, handle func(pass string),
) *Password {
	pass := w.Editor().Mask('•').SingleLine().Submit(true)
	passInput := w.TextInput(pass, hint).Color(borderColorUnfocused)
	p := &Password{
		Window:               w,
		unhideClickable:      w.Clickable(),
		copyClickable:        w.Clickable(),
		pasteClickable:       w.Clickable(),
		pass:                 pass,
		passInput:            passInput,
		borderColorUnfocused: borderColorUnfocused,
		borderColorFocused:   borderColorFocused,
		borderColor:          borderColorUnfocused,
		backgroundColor:      backgroundColor,
		handle:               handle,
		password:             password,
	}
	p.copyButton = w.IconButton(p.copyClickable)
	p.pasteButton = w.IconButton(p.pasteClickable)
	p.unhideButton = w.IconButton(p.unhideClickable).Background("").
		Icon(w.Icon().Color(p.borderColor).Src(&icons2.ActionVisibility))
	p.showClickableFn = func(col string) {
		Debug("show clickable clickeded")
		p.hide = !p.hide
		// if !p.hide {
		// 	p.unhideButton = p.unhideButton.
		// 		Color("Primary").
		// 		Icon(
		// 			Theme.Icon().
		// 				Color(col).
		// 				Src(&icons2.ActionVisibility),
		// 		)
		// 	// p.pass.Mask('•')
		// 	p.passInput.Color(col)
		// } else {
		// 	p.unhideButton = p.unhideButton.
		// 		Color("DocText").
		// 		Icon(
		// 			Theme.Icon().
		// 				Color(p.borderColor).
		// 				Src(&icons2.ActionVisibilityOff),
		// 		)
		// 	// p.pass.Mask(0)
		// 	p.passInput.Color(col)
		// }
		p.pass.Focus()
	}
	copyClickableFn := func() {
		go clipboard.Set(p.pass.Text())
		p.pass.Focus()
	}
	pasteClickableFn := func() {
		go func() {
			txt := p.pass.Text()
			txt = txt[:p.pass.Caret.Col] + clipboard.Get() + txt[p.pass.Caret.Col:]
			p.pass.SetText(txt)
		}()
		p.pass.Focus()
	}
	p.copyClickable.SetClick(copyClickableFn)
	p.pasteClickable.SetClick(pasteClickableFn)
	p.unhideButton.
		// Color("Primary").
		Icon(
			w.Icon().
				Color(p.borderColor).
				Src(&icons2.ActionVisibility),
		)
	p.pass.Mask('•')
	p.pass.SetFocus(
		func(is bool) {
			if is {
				p.borderColor = p.borderColorFocused
			} else {
				p.borderColor = p.borderColorUnfocused
				p.hide = true
			}
		},
	)
	p.passInput.editor.Mask('•')
	p.hide = true
	p.passInput.Color(p.borderColor)
	p.pass.SetText(*p.password).Mask('•').SetSubmit(
		func(txt string) {
			// if !p.hide {
			// 	p.showClickableFn(p.borderColor)
			// }
			// p.showClickableFn(p.borderColor)
			go func() {
				p.handle(txt)
			}()
		},
	).SetChange(
		func(txt string) {
			// send keystrokes to the NSA
		},
	)
	return p
}

func (p *Password) Fn(gtx l.Context) l.Dimensions {
	// gtx.Constraints.Max.X = int(p.TextSize.Scale(float32(p.size)).V)
	// gtx.Constraints.Min.X = 0
	// cs := gtx.Constraints
	// width := int(p.Theme.TextSize.Scale(p.size).V)
	// gtx.Constraints.Max.X, gtx.Constraints.Min.X = width, width
	return func(gtx l.Context) l.Dimensions {
		p.passInput.Color(p.borderColor).Font("go regular")
		p.unhideButton.Color(p.borderColor)
		p.unhideClickable.SetClick(func() { p.showClickableFn(p.borderColor) })
		visIcon := &icons2.ActionVisibility
		if p.hide {
			p.pass.Mask('•')
		} else {
			visIcon = &icons2.ActionVisibilityOff
			p.pass.Mask(0)
		}
		
		return p.Fill(
			p.backgroundColor, l.Center, 0, 0,
			p.Inset(
				0.25,
				p.Flex().
					Flexed(
						1,
						p.Inset(0.25, p.passInput.Color(p.borderColor).HintColor(p.borderColorUnfocused).Fn).Fn,
					).
					Rigid(
						p.copyButton.
							Background("").
							Icon(p.Icon().Color(p.borderColor).Scale(Scales["H6"]).Src(&icons2.ContentContentCopy)).
							ButtonInset(0.25).
							Fn,
					).
					Rigid(
						p.pasteButton.
							Background("").
							Icon(p.Icon().Color(p.borderColor).Scale(Scales["H6"]).Src(&icons2.ContentContentPaste)).
							ButtonInset(0.25).
							Fn,
					).
					Rigid(
						p.unhideButton.
							Background("").
							Icon(p.Icon().Color(p.borderColor).Src(visIcon)).Fn,
					).
					Fn,
			).Fn,
		).Fn(gtx)
	}(gtx)
}

func (p *Password) GetPassword() string {
	return p.passInput.editor.Text()
}

func (p *Password) Wipe() {
	p.passInput.editor.rr.Zero()
	p.passInput.editor.SetText("")
}

func (p *Password) Focus() {
	p.passInput.editor.Focus()
}

func (p *Password) Blur() {
	p.passInput.editor.focused = false
}

func (p *Password) Hide() {
	p.passInput.editor.Mask('*')
}

func (p *Password) Show() {
	p.passInput.editor.Mask(0)
}
