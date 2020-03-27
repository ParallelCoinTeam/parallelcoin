package main

import (
	"github.com/p9c/logi"
	log "github.com/p9c/pod/pkg/logi"
)

func main() {

}

// Serve is the API for a serve endpoint (program producing logs)
type Serve struct{}

func (s *Serve) Run(_ *struct{}, ack *bool) (err error) {
	return
}

func (s *Serve) Stop(_ *struct{}, ack *bool) (err error) {
	return
}

func (s *Serve) Log(ent log.Entry) {}

// Consume is the API for a consume endpoint (program consuming logs)
type Consume struct{}

func (c *Consume) Run() (ack bool) {
	return
}

func (c *Consume) Stop() (ack bool) {
	return
}

func (c *Consume) Log(ent *logi.Entry, ack *bool) (err error) {
	return
}