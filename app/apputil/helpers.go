package apputil

import (
	"os"
	"time"

	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
)

// NewCommand returns a cli.Command
func NewCommand(
	name string,
	usage string,
	action interface{},
	subcommands cli.Commands,
	aliases ...string,
) cli.Command {
	return cli.Command{
		Name:        name,
		Aliases:     aliases,
		Usage:       usage,
		Action:      action,
		Subcommands: subcommands,
	}
}

// SubCommands returns a slice of cli.Command
func SubCommands(sc ...cli.Command) []cli.Command {
	return append([]cli.Command{}, sc...)
}

// String returns an altsrc.StringFlag
func String(name, usage, value string, dest *string) *altsrc.StringFlag {
	return altsrc.NewStringFlag(cli.StringFlag{
		Name:        name,
		Usage:       usage,
		Value:       value,
		Destination: dest,
	})
}

// BoolTrue returns a CliBoolFlag that defaults to true
func BoolTrue(name, usage string, dest *bool) *altsrc.BoolTFlag {
	return altsrc.NewBoolTFlag(cli.BoolTFlag{
		Name:        name,
		Usage:       usage,
		Destination: dest,
	})
}

// Bool returns an altsrc.BoolFlag
func Bool(name, usage string, dest *bool) *altsrc.BoolFlag {
	return altsrc.NewBoolFlag(cli.BoolFlag{
		Name:        name,
		Usage:       usage,
		Destination: dest,
	})
}

// Join joins together a path and filename
func Join(path, filename string) string {
	return path + string(os.PathSeparator) + filename
}

// StringSlice returns and altsrc.StringSliceFlag
func StringSlice(name, usage string, value *cli.StringSlice) *altsrc.StringSliceFlag {
	return altsrc.NewStringSliceFlag(cli.StringSliceFlag{
		Name:  name,
		Usage: usage,
		Value: value,
	})
}

// Int returns an altsrc.IntFlag
func Int(name, usage string, value int, dest *int) *altsrc.IntFlag {
	return altsrc.NewIntFlag(cli.IntFlag{
		Name:        name,
		Value:       value,
		Usage:       usage,
		Destination: dest,
	})
}

// Duration returns an altsrc.DurationFlag
func Duration(name, usage string, value time.Duration, dest *time.Duration) *altsrc.DurationFlag {
	return altsrc.NewDurationFlag(cli.DurationFlag{
		Name:        name,
		Value:       value,
		Usage:       usage,
		Destination: dest,
	})
}

// Float64 returns an altsrc.Float64Flag
func Float64(name, usage string, value float64, dest *float64) *altsrc.Float64Flag {
	return altsrc.NewFloat64Flag(cli.Float64Flag{
		Name:        name,
		Value:       value,
		Usage:       usage,
		Destination: dest,
	})
}
