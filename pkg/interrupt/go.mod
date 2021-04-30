module github.com/p9c/interrupt

go 1.16

require (
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
github.com/p9c/interrupt
	github.com/p9c/log v0.0.4
	github.com/p9c/qu v0.0.3
	github.com/stretchr/testify v1.6.1 // indirect
	go.uber.org/atomic v1.7.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1
)
replace (
	github.com/p9c/log => ../log
	github.com/p9c/qu => ../qu
	github.com/p9c/interrupt => ../interrupt
)
