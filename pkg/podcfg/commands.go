package podcfg

// Commands are a slice of podcfg.Command entries
type Commands []Command

// Command is a specification for a command and can include any number of subcommands
type Command struct {
	Name        string
	Description string
	Entrypoint  func(c *Config) error
	Commands    Commands
}

var tabs = "\t\t\t\t\t"

// Find the Command you are looking for. Note that the namespace is assumed to be flat, no duplicated names on different
// levels, as it returns on the first one it finds, which goes depth-first recursive
func (c Commands) Find(name string, hereDepth, hereDist int) (found bool, depth, dist int, cm *Command, e error) {
	if c == nil {
		dist = hereDist
		depth = hereDepth
		return
	}
	if hereDist == 0 {
		D.Ln("searching for command:", name)
	}
	depth = hereDepth + 1
	T.Ln(tabs[:depth]+"->", depth)
	dist = hereDist
	for i := range c {
		T.Ln(tabs[:depth]+"walking", c[i].Name, depth, dist)
		if c[i].Name == name {
			// depth++
			T.Ln(tabs[:depth]+"found", name, "at depth", depth, "distance", dist)
			found = true
			cm = &c[i]
			e = nil
			return
		} else {
			dist++
		}
		if found, depth, dist, cm, e = c[i].Commands.Find(name, depth, dist); E.Chk(e) {
			T.Ln(tabs[:depth]+"error", c[i].Name)
			// depth--
			return
		}
		if found {
			return
		}
	}
	T.Ln(tabs[:hereDepth]+"<-", hereDepth)
	if hereDepth == 0 {
		D.Ln("search text", name, "not found")
	}
	depth--
	return
}
