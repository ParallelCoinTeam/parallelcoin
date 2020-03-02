package model

type DuoUIdialog struct {
	Show   bool
	Ok     func()
	Close  func()
	Cancel func()
	Title  string
	Text   string
}
