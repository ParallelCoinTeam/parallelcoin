package res

import (
	"fmt"
	"gioui.org/layout"
)

type Responsive struct {
	Mode string
	Mod  map[string]interface{}
}

func Resposnsivity(x, y int) *Responsive {
	//fmt.Println("ResposnsivityResposnsivityResposnsivity", x)

	switch w := x; {
	//case w > 2560:
	//	g.UI.Res.Mode = "desktopXL"
	//case w > 1920:
	//	g.UI.Res.Mode = "desktopL"
	//case w > 1680:
	//	g.UI.Res.Mode = "desktopM"
	//case w > 1440:z
	//	g.UI.Res.Mode = "desktopS"
	case w > 1024:
		fmt.Println("Screen 1024")
		return Desktop()

		//r.Mod["Screen"] = "max(inset(0dp90dp50dp10dp,_))"
		//r.Mod["Main"] = "max(hflex(start,r(_),f(1,_)))"
		//r.Mod["Nav"] = "vflex(start,r(_),f(1,_))"
		//r.Mod["NavItemsAxis"] = layout.Vertical
		//r.Mod["Content"] = "max(vflex(start,r(_),f(1,_),r(_)))"
		//r.Mod["TwoEqual"] = "max(hflex(start,f(0.5,_),f(0.5,_)))"
	//case w > 960:
	//	fmt.Println("Screen 900")
	//	r.Mode = "tablet"
	//
	//case w > 600:
	//	fmt.Println("Screen 600")
	//r.Mod["TwoEqual"] = "vflex(start,f(0.5,_),f(0.5,_))"
	//case w > 1136:
	//	g.UI.Res.Mode = "tabletS"
	//case w > 1024:
	//	g.UI.Res.Mode = "tabletM"
	//case w > 960:
	//	g.UI.Res.Mode = "mobileXXXL"
	//case w > 800:
	//	g.UI.Res.Mode = "mobileXXL"
	//case w > 640:
	//	g.UI.Res.Mode = "mobileXL"
	//case w > 480:
	//	g.UI.Res.Mode = "mobileL"
	//case w > 320:
	//	g.UI.Res.Mode = "mobileM"
	//case w > 240:
	//	g.UI.Res.Mode = "mobileS"
	default:
		fmt.Println("Screen 00")
		return Mobile()
	}
}

//func Resposnsivity(r *Responsivity) {
//
//	return
//}

func Desktop() *Responsive {
	return &Responsive{
		Mode: "desktop",
		Mod: map[string]interface{}{
			"Screen":    "max(inset(80dp0dp80dp0dp,_))",
			"Container": "max(hflex(start,r(_),f(1,_)))",
			"Main":      "max(vflex(start,f(1,_),r(_)))",

			"Page":     "max(vflex(start,r(_),f(1,_),r(_)))",
			"TwoEqual": "hflex(start,f(0.5,inset(0dp,_)),f(0.5,inset(0dp,_))))",

			"Nav":             "vflex(middle,r(_),f(1,inset(0dp0dp0dp0dp,_)))",
			"NavSize":         128,
			"NavIconSize":     128,
			"NavIconAndLabel": "hflex(start, r(_),f(1,_))",
			//"Logo":            d.UI.Theme.Icons["Logo"],
			"NavItemsAxis": layout.Vertical,
		},
	}
}

func Mobile() *Responsive {
	return &Responsive{
		Mode: "mobile",
		Mod: map[string]interface{}{
			"Screen":    "max(inset(80dp0dp80dp0dp,_))",
			"Container": "max(vflex(start,r(_),f(1,_)))",
			"Main":      "max(vflex(start,f(1,_),r(_)))",

			"Page":     "max(vflex(start,r(_),f(1,_),r(_)))",
			"TwoEqual": "vflex(start,r(inset(0dp,_)),f(1,inset(0dp,_))))",

			"Nav":             "hflex(middle,r(_),f(1,inset(0dp0dp0dp0dp,_)))",
			"NavSize":         128,
			"NavIconSize":     128,
			"NavIconAndLabel": "hflex(r(_),r(_))",
			//"Logo":            d.UI.Theme.Icons["Logo"],
			"NavItemsAxis": layout.Horizontal,
		},
	}
}
