package mod

//  EJS Layout Model
type EJSlayout struct {
	AllowDragging     bool       `json:"allowDragging"`
	AllowFloating     bool       `json:"allowFloating"`
	AllowResizing     bool       `json:"allowResizing"`
	CellAspectRatio   string     `json:"cellAspectratio"`
	CellSpacing       []int      `json:"cellSpacing"`
	Columns           int        `json:"columns"`
	DraggableHandle   string     `json:"draggaBlehandle"`
	EnablePersistence bool       `json:"enablePersistence"`
	EnableRtl         bool       `json:"enableRtl"`
	MediaQuery        string     `json:"mediaQuery"`
	Panels            []EJSpanel `json:"panels"`
	ResizableHandles  []string   `json:"resizableHandles"`
	ShowGridLines     bool       `json:"showGridLines"`
}

//  EJS Panel Model
type EJSpanel struct {
	ID       string `json:"id"`
	Col      int    `json:"col"`
	Content  string `json:"content"`
	CssClass string `json:"cssClass"`
	Enabled  bool   `json:"enabled"`
	Header   string `json:"header"`
	MaxSizeX int    `json:"maxSizeX"`
	MaxSizeY int    `json:"maxSizeY"`
	MinSizeX int    `json:"minSizeX"`
	MinSizeY int    `json:"minSizeY"`
	Row      int    `json:"row"`
	SizeX    int    `json:"sizeX"`
	SizeY    int    `json:"sizeY"`
	ZIndex   int    `json:"zIndex"`
}
