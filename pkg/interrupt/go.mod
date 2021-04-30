module github.com/p9c/interrupt

go 1.16

require (
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/p9c/log v0.0.12
	github.com/p9c/qu v0.0.12
	github.com/xanzy/ssh-agent v0.2.1 // indirect
	go.uber.org/atomic v1.7.0
	gopkg.in/src-d/go-billy.v4 v4.3.2 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1
)

replace (
	github.com/p9c/log => ../log
	github.com/p9c/qu => ../qu
)
