package model

type
	DuoUIdialog struct {
		Show   bool
		Ok     func()
		Cancel func()
		Title  string
		Text   string
	}
