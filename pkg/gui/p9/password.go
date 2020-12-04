package p9

import (
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"

	l "gioui.org/layout"

	"github.com/p9c/pod/pkg/gui/clipboard"
)

type Password struct {
	*Theme
	pass                 *Editor
	passInput            *TextInput
	unhideClickable      *Clickable
	unhideButton         *IconButton
	copyClickable        *Clickable
	copyButton           *IconButton
	pasteClickable       *Clickable
	pasteButton          *IconButton
	hide                 bool
	size                 float32
	borderColor          string
	borderColorUnfocused string
	borderColorFocused   string
	focused              bool
	showClickableFn      func(col string)
	password             *string
	handle               func(pass string)
}

func (th *Theme) Password(hint string, password *string, borderColorFocused, borderColorUnfocused string, size float32, handle func(pass string)) *Password {
	pass := th.Editor().Mask('•').SingleLine().Submit(true)
	passInput := th.TextInput(pass, hint).Color(borderColorUnfocused)
	p := &Password{
		Theme:                th,
		unhideClickable:      th.Clickable(),
		copyClickable:        th.Clickable(),
		pasteClickable:       th.Clickable(),
		pass:                 pass,
		passInput:            passInput,
		size:                 size,
		borderColorUnfocused: borderColorUnfocused,
		borderColorFocused:   borderColorFocused,
		borderColor:          borderColorUnfocused,
		handle:               handle,
		password:             password,
	}
	p.copyButton = th.IconButton(p.copyClickable)
	p.pasteButton = th.IconButton(p.pasteClickable)
	p.unhideButton = th.IconButton(p.unhideClickable).
		Background("Transparent").
		Icon(th.Icon().Color(p.borderColor).Src(&icons2.ActionVisibility))
	p.showClickableFn = func(col string) {
		p.hide = !p.hide
		if !p.hide {
			p.unhideButton.
				// Color("Primary").
				Icon(
					th.Icon().
						Color(col).
						Src(&icons2.ActionVisibility))
			p.pass.Mask('•')
			p.passInput.Color(col)
		} else {
			p.unhideButton.
				// Color("DocText").
				Icon(
					th.Icon().
						Color(p.borderColor).
						Src(&icons2.ActionVisibilityOff),
				)
			p.pass.Mask(0)
			p.passInput.Color(col)
		}
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
			th.Icon().
				Color(p.borderColor).
				Src(&icons2.ActionVisibility),
		)
	p.pass.Mask('•')
	p.pass.SetFocus(func(is bool) {
		if is {
			p.borderColor = p.borderColorFocused
		} else {
			p.borderColor = p.borderColorUnfocused
			p.hide = true
		}
	})
	p.passInput.Color(p.borderColor)
	p.pass.SetText(*p.password).Mask('•').SetSubmit(func(txt string) {
		// if !p.hide {
		// 	p.showClickableFn(p.borderColor)
		// }
		// p.showClickableFn(p.borderColor)
		go func() {
			p.handle(txt)
		}()
	}).SetChange(func(txt string) {
		// send keystrokes to the NSA
	})
	return p
}

func (p *Password) Fn(gtx l.Context) l.Dimensions {
	// gtx.Constraints.Max.X = int(p.TextSize.Scale(float32(p.size)).V)
	// gtx.Constraints.Min.X = 0
	// cs := gtx.Constraints
	width := int(p.Theme.TextSize.Scale(p.size).V)
	gtx.Constraints.Max.X, gtx.Constraints.Min.X = width, width
	return func(gtx l.Context) l.Dimensions {
		p.passInput.Color(p.borderColor)
		p.unhideButton.Color(p.borderColor)
		p.unhideClickable.SetClick(func() { p.showClickableFn(p.borderColor) })
		if p.hide {
			p.pass.Mask('•')
		} else {
			p.pass.Mask(0)
		}

		return p.Border().Color(p.borderColor).Embed(
			p.Flex().
				Flexed(1,
					p.Inset(0.25, p.passInput.Color(p.borderColor).Fn).Fn,
				).
				Rigid(
					p.copyButton.
						Background("").
						Icon(p.Icon().Color(p.borderColor).Scale(Scales["H6"]).Src(&icons2.ContentContentCopy)).
						Inset(0.25).
						Fn,
				).
				Rigid(
					p.pasteButton.
						Background("").
						Icon(p.Icon().Color(p.borderColor).Scale(Scales["H6"]).Src(&icons2.ContentContentPaste)).
						Inset(0.25).
						Fn,
				).
				Rigid(
					p.unhideButton.
						Background("Transparent").
						Icon(p.Icon().Color(p.borderColor).Src(&icons2.ActionVisibility)).Fn,
				).
				Fn,
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
