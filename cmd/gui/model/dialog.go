package model

type DuoUIdialog struct {
	Show        bool
	Ok          func()
	Close       func()
	Cancel      func()
	CustomField func()
	Title       string
	Text        string
}

