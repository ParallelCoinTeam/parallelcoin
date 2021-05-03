module github.com/p9c/gel

go 1.16

require (
	github.com/BurntSushi/xgb v0.0.0-20210121224620-deaf085860bc
	github.com/atotto/clipboard v0.1.4
	github.com/p9c/gio v0.0.5
	github.com/p9c/interrupt v0.0.8
	github.com/p9c/log v0.0.12
	github.com/p9c/opts v0.0.13
	github.com/p9c/qu v0.0.12
	go.uber.org/atomic v1.7.0
	golang.org/x/exp v0.0.0-20210417010653-0739314eea07
	golang.org/x/image v0.0.0-20210220032944-ac19c3e999fb
	gopkg.in/src-d/go-git.v4 v4.13.1
)

replace (
	github.com/p9c/gio => ./gio
	github.com/p9c/log => ../log
	github.com/p9c/opts => ../opts
	github.com/p9c/qu => ../qu
)
