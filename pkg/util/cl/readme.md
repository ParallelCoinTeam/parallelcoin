# clog

Clog is a logger that works by receiving values, including closures, on channels and triggering an output through a configurable output function. The log level can be changed on the fly and immediately takes effect. There is colour available for tags if desired.

## Documentation

See here at our private [godoc](http://89.40.12.55:8008/pkg/git.parallelcoin.io/dev/pod/pkg/util/cl/) (or go get it and run your own :) )

## Automatic

As soon as you import it into your application it is instantly activated by an automatic init routine, and is available package-wide, enabling easy per-subsystem logging based on the application's filesystem hierarchy, as well as within packages, it is a one liner to declare a new subsystem.

## Channel-based

Because the loggers are channel based, the overhead of logging function calls, which can be quite extensive, is completely eliminated from the main goroutines the logging calls are located.

Each log call only costs the processing time for loading a channel and notifying the main goroutine manager that the channel now has a new thing in it, which will run in a separate goroutine from where the log strings were sent.

To further reduce interruptions inside intensive, long running loops, there is a 'closure' type which can be used to embed more expensive queries for internal data that will only incur a processing cost if the log level is enabled.

Thus, in the background thread running the subsystem logger, when the log level enabled is lower than the specific entry in question, the cost will be in receiving the log string or pointer to the closure, and then a comparison and if it aborts, the logger returns to waiting for the next log entry to arrive on its given channel.

Each subsystem spawns its own goroutine, and disabled ones just discard, so even within the subsystem, actually processing enabled log levels are in separate threads as well.

## Subsystems

One simply calls `clog.NewSubSystem` with a name and level and you get back a struct containing named channels that sit waiting to prepend your prescribed name string on whatever you put in the channel, and then forward it to the appropriate level of the root clog logger (accessible through `clog.L`).

It is recommended to make all of the names right-padded with spaces in order to keep the start of log entries in a consistent position, this is not handled automatically.

### Example launch and log

```go
ss := clog.NewSubSystem("TEST", clog.Ndbg)
ss.Info <- "this will print an info message"
```

Output will be like this:

```
2019-01-15 11:59:58.155324 UTC [INF] this will print an info message
```

Note also that by default, the log level identifiers are colour coded to help more easily distinguish between log types.
