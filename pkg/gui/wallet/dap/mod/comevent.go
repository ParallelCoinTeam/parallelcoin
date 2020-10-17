package mod

type Commands struct {
	Events  chan CommandsEvent
	History []Command
}

type Command struct {
	Command func()
}

type CommandsEvent interface {
	isCommandsEvent()
}

type CommandEvent struct {
	Command Command
}

type ErrorEvent struct {
	Err error
}

func (d *Commands) Run() {
	d.Events <- CommandEvent{
		Command: Command{},
	}
}

func (e CommandEvent) isCommandsEvent() {}
func (e ErrorEvent) isCommandsEvent()   {}
