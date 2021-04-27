package cmds

// Commands are a slice of Command entries
type Commands []Command

// Command is a specification for a command and can include any number of subcommands
type Command struct {
	Name          string
	Description   string
	Documentation string
	Entrypoint    func(c interface{}) error
	Commands      Commands
	Colorizer     func(a ...interface{}) string
	AppText       string
	Parent        *Command
}

func (c Commands) PopulateParents(parent *Command) {
	if parent != nil {
		T.Ln("backlinking children of", parent.Name)
	}
	for i := range c {
		c[i].Parent = parent
		c[i].Commands.PopulateParents(&c[i])
	}
}

// GetAllCommands returns all of the available command names
func (c Commands) GetAllCommands() (o []string) {
	c.ForEach(func(cm Command) bool {
		o = append(o, cm.Name)
		o = append(o, cm.Commands.GetAllCommands()...)
		return true
	}, 0, 0,
	)
	return
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

func (c Commands) ForEach(fn func(Command) bool, hereDepth, hereDist int) (ret bool, depth, dist int, e error) {
	if c == nil {
		dist = hereDist
		depth = hereDepth
		return
	}
	depth = hereDepth + 1
	T.Ln(tabs[:depth]+"->", depth)
	dist = hereDist
	for i := range c {
		T.Ln(tabs[:depth]+"walking", c[i].Name, depth, dist)
		if !fn(c[i]) {
			// if the closure returns false break out of the loop
			return
		}
	}
	T.Ln(tabs[:hereDepth]+"<-", hereDepth)
	depth--
	return
}
