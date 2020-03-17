package rcd

type DuoUIcommands struct {
	Events  chan DuoUIcommandsEvent
	History []DuoUIcommand
}

type DuoUIcommand struct {
	Command func()
}

type DuoUIcommandsEvent interface {
	isCommandsEvent()
}

type DuoUIcommandEvent struct {
	Command DuoUIcommand
}

type ErrorEvent struct {
	Err error
}

func (d *DuoUIcommands) Run() {
	d.Events <- DuoUIcommandEvent{
		Command: DuoUIcommand{},
	}
}

func (e DuoUIcommandEvent) isCommandsEvent() {}
func (e ErrorEvent) isCommandsEvent()        {}
