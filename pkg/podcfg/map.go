package podcfg

// MapGen takes a Config and gives it a map
func MapGen(in *Config) (c *Config) {
	c = in
	c.Map = make(map[string]interface{})
	c.ForEach(
		func(ifc interface{}) bool {
			switch ii := ifc.(type) {
			case *Bool:
				im := ii.Metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Option]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Option] = ii
			case *Strings:
				im := ii.Metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Option]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Option] = ii
			case *Float:
				im := ii.Metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Option]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Option] = ii
			case *Int:
				im := ii.Metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Option]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Option] = ii
			case *String:
				im := ii.Metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Option]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Option] = ii
			case *Duration:
				im := ii.Metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Option]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Option] = ii
			default:
			}
			return true
		},
	)
	return
}
