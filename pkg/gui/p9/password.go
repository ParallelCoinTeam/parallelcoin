package p9

import (
	l "gioui.org/layout"
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"
)

type Password struct {
	*Theme
	pass            *Editor
	passInput       *TextInput
	unhideClickable *Clickable
	unhideButton    *IconButton
	GetPassword     func() string
	hide            bool
}

func (th *Theme) Password(password *string, handle func(pass string)) *Password {
	pass := th.Editor().Mask('•').SingleLine(true).Submit(true)
	passInput := th.SimpleInput(pass).Color("DocText")
	p := &Password{
		Theme:           th,
		unhideButton:    nil,
		unhideClickable: th.Clickable(),
		pass:            pass,
		passInput:       passInput,
	}
	p.GetPassword = func() string {
		return p.pass.Text()
	}
	p.unhideButton = th.IconButton(p.unhideClickable).
		Background("").
		Icon(th.Icon().Color("Primary").Src(icons2.ActionVisibility))
	showClickableFn := func() {
		p.hide = !p.hide
		if !p.hide {
			p.unhideButton.
				// Color("Primary").
				Icon(
					th.Icon().
						Color("Primary").
						Src(icons2.ActionVisibility))
			p.pass.Mask('•')
			p.passInput.Color("Primary")
		} else {
			p.unhideButton.
				// Color("DocText").
				Icon(
					th.Icon().
						Color("DocText").
						Src(icons2.ActionVisibilityOff),
				)
			p.pass.Mask(0)
			p.passInput.Color("DocText")
		}
	}
	p.unhideButton.
		// Color("Primary").
		Icon(
			th.Icon().
				Color("Primary").
				Src(icons2.ActionVisibility),
		)
	p.pass.Mask('•')
	p.passInput.Color("Primary")
	p.unhideClickable.SetClick(showClickableFn)
	p.pass.SetText(*password).Mask('•').SetSubmit(func(txt string) {
		if !p.hide {
			showClickableFn()
		}
		showClickableFn()
		go func() {
			handle(txt)
		}()
	}).SetChange(func(txt string) {
		// send keystrokes to the NSA
	})

	return p
}

func (p *Password) Fn(gtx l.Context) l.Dimensions {
	return p.Border().Embed(
		p.Flex().Flexed(1,
			p.Inset(0.25, p.passInput.Fn).Fn,
		).Rigid(
			p.unhideButton.Fn,
		).Fn,
	).Fn(gtx)
}
