// SPDX-License-Identifier: Unlicense OR MIT

// Package profiles provides access to rendering
// profiles.
package profile

import (
	"github.com/p9c/pod/pkg/gio/internal/opconst"
	"github.com/p9c/pod/pkg/gio/io/event"
	"github.com/p9c/pod/pkg/gio/op"
)

// Op registers a handler for receiving
// Events.
type Op struct {
	Key event.Key
}

// Event contains profile data from a single
// rendered frame.
type Event struct {
	// Timings. Very likely to change.
	Timings string
}

func (p Op) Add(o *op.Ops) {
	data := o.Write(opconst.TypeProfileLen, p.Key)
	data[0] = byte(opconst.TypeProfile)
}

func (p Event) ImplementsEvent() {}
