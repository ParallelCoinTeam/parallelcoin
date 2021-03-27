package podcfg

// MapGen takes a Config and gives it a map
func MapGen(in *Config) (c *Config) {
	c = in
	c.Map = make(map[string]interface{})
	c.ForEach(
		func(ifc interface{}) bool {
			switch ii := ifc.(type) {
			case *Bool:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Name] = ii
			case *Strings:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Name] = ii
			case *Float:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Name] = ii
			case *Int:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Name] = ii
			case *String:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Name] = ii
			case *Duration:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Name] = ii
			default:
			}
			return true
		},
	)
	return
}
