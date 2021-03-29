package podcfg

import (
	"fmt"
)

func (c *Config) getAllOptionStrings() (s map[string][]string, e error) {
	s = make(map[string][]string)
	if c.ForEach(func(ifc Option) bool {
		md := ifc.GetMetadata()
		if _, ok := s[ifc.Name()]; ok {
			e = fmt.Errorf("conflicting option names: %v %v", ifc.GetAllOptionStrings(), s[ifc.Name()])
			return false
		}
		s[ifc.Name()] = md.GetAllOptionStrings()
		return true
	},
	) {
	}
	s["commandslist"] = c.Commands.GetAllCommands()
	// I.S(s["commandslist"])
	return
}

func findConflictingItems(valOpts map[string][]string) (o []string, e error) {
	var ss, ls string
	for i := range valOpts {
		for j := range valOpts {
			// W.Ln(s[i], s[j], i==j, s[i]==s[j])
			if i == j {
				continue
			}
			a := valOpts[i]
			b := valOpts[j]
			for ii := range a {
				for jj := range b {
					if ii == jj {
						continue
					}
					// W.Ln(i == j, s[i] == s[j])
					// I.Ln(s[i], s[j])
					ss, ls = shortestString(a[ii], b[jj])
					// I.Ln("these should not be the same string", ss, ls)
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
