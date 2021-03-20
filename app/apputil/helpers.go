package apputil

import (
	"os"
	"time"
	
	"github.com/urfave/cli"
)

// Command returns a cli.Command
func Command(
	name string,
	usage string,
	action interface{},
	subcommands cli.Commands,
	flags []cli.Flag,
	aliases ...string,
) cli.Command {
	return cli.Command{
		Name:        name,
		Aliases:     aliases,
		Usage:       usage,
		Action:      action,
		Subcommands: subcommands,
		Flags:       flags,
	}
}

// SubCommands returns a slice of cli.Command
func SubCommands(sc ...cli.Command) []cli.Command {
	return append([]cli.Command{}, sc...)
}

// Lang returns an cli.StringFlag
func Lang(name, usage, value string, dest *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        name,
		Usage:       usage,
		Value:       value,
		Destination: dest,
	}
}

// String returns an cli.StringFlag
func String(name, usage, value string, dest *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        name,
		Usage:       usage,
		Value:       value,
		Destination: dest,
	}
}

// BoolTrue returns a CliBoolFlag that defaults to true
func BoolTrue(name, usage string, dest *bool) *cli.BoolTFlag {
	return &cli.BoolTFlag{
		Name:        name,
		Usage:       usage,
		Destination: dest,
	}
}

// Bool returns an cli.BoolFlag
func Bool(name, usage string, dest *bool) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        name,
		Usage:       usage,
		Destination: dest,
	}
}

// Join joins together a path and filename
func Join(path, filename string) string {
	return path + string(os.PathSeparator) + filename
}

// StringSlice returns and cli.StringSliceFlag
func StringSlice(name, usage string, value *cli.StringSlice) *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:  name,
		Usage: usage,
		Value: value,
	}
}

// Int returns an cli.IntFlag
func Int(name, usage string, value int, dest *int) *cli.IntFlag {
	return &cli.IntFlag{
		Name:        name,
		Value:       value,
		Usage:       usage,
		Destination: dest,
	}
}

// Duration returns an cli.DurationFlag
func Duration(name, usage string, value time.Duration, dest *time.Duration) *cli.DurationFlag {
	return &cli.DurationFlag{
		Name:        name,
		Value:       value,
		Usage:       usage,
		Destination: dest,
	}
}

// Float64 returns an cli.Float64Flag
func Float64(name, usage string, value float64, dest *float64) *cli.Float64Flag {
	return &cli.Float64Flag{
		Name:        name,
		Value:       value,
		Usage:       usage,
		Destination: dest,
	}
}
