// Package cl is clog, the channel logger
//
// What's this for
//
// Many processes involve very tight loops that number in the hundreds of nanoseconds per iteration, and logging systems can contribute a significant additional time to this.
//
// This library provides a logging subsystem that works by pushing logging data into channels and allowing the cost of logging to be minimised especially in tight loop situations. To this end also there is a closure channel type that lets you defer the query of data for a log indefinitely if the log level is currenntly inactive.
//
// The main benefit of using channels to coordinate logging is that it allows logging to occupy a separate thread to execution, meaning issues involving blocking on the processing but especially output to pipes and tty devices never affects main loops directly.
package cl
