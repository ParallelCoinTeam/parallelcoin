package main

import (
	"fmt"
	"github.com/gookit/color"
	"os"
	"strings"
)

var flags = []string{"wc", "discover", "l", "g", "D", "solo", "G", "norpc", "G", "autoports"}
var commands = []string{"node", "gui", "wallet", "ctl"}

func main() {
	// for i := range os.Args {
	// 	fmt.Printf("'%s' ", os.Args[i])
	// }
	// fmt.Println()
	// first find the first command and split the args
	cmdPos := -1
next:
	for i := range os.Args {
		if i == 0 {
			continue
		}
		if cmdPos < 0 {
			for j := range flags {
				if strings.HasPrefix(os.Args[i], flags[j]) {
					fmt.Println("recognised flag", flags[j], os.Args[i])
					continue next
				}
			}
			for j := range commands {
				if commands[j] == os.Args[i] {
					fmt.Println("top level command found", os.Args[i])
					cmdPos = i
					continue next
				}
			}
			var o string
			for j := range os.Args {
				if j == i {
					o += color.Red.Sprint(os.Args[j] + " ")
				} else {
					o += os.Args[j] + " "
				}
			}
			fmt.Printf("ABORT: unrecognised flag '%s'\n", o)
			break next
		} else {
			for j := range commands {
				if commands[j] == os.Args[i] {
					fmt.Println("added subcommand", os.Args[i])
					continue next
				}
			}
			var o string
			for j := range os.Args {
				if j == i {
					o += color.Red.Sprint(os.Args[j])
				} else {
					o += os.Args[j]
				}
				if j < len(os.Args)-1 {
					o += " "
				}
			}
			fmt.Printf("ABORT: unrecognised subcommand '%s'\n", o)
			break next
		}
	}
	// for i := range os.Args {
	// 	if cmdPos == i {
	// 		fmt.Println("\ncommands")
	// 	}
	// 	fmt.Printf("'%s' ", os.Args[i])
	// }
	// fmt.Println()
}
