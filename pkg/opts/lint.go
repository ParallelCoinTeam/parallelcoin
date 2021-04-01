package opts

import (
	"fmt"
	"github.com/p9c/opts/opt"
)

func getAllOptionStrings(c *Config) (s map[string][]string, e error) {
	s = make(map[string][]string)
	if c.ForEach(func(ifc opt.Option) bool {
		md := ifc.GetMetadata()
		if _, ok := s[ifc.Name()]; ok {
			e = fmt.Errorf("conflicting opt names: %v %v", ifc.GetAllOptionStrings(), s[ifc.Name()])
			return false
		}
		s[ifc.Name()] = md.GetAllOptionStrings()
		return true
	},
	) {
	}
	s["commandslist"] = c.Commands.GetAllCommands()
	return
}

func findConflictingItems(valOpts map[string][]string) (o []string, e error) {
	var ss, ls string
	for i := range valOpts {
		for j := range valOpts {
			if i == j {
				continue
			}
			a := valOpts[i]
			b := valOpts[j]
			for ii := range a {
				for jj := range b {
					ss, ls = shortestString(a[ii], b[jj])
					if ss == ls[:len(ss)] {
						E.F("conflict between %s and %s, ", ss, ls)
						o = append(o, ss, ls)
					}
				}
			}
		}
	}
	if len(o) > 0 {
		panic(fmt.Sprintf("conflicts found: %v", o))
	}
	return
}

func shortestString(a, b string) (s, l string) {
	switch {
	case len(a) > len(b):
		s, l = b, a
	default:
		s, l = a, b
	}
	return
}
