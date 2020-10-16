package theme

import (
	"gioui.org/widget/material"
)

func H1(th *Theme, txt string) material.LabelStyle {
	return material.Label(th.T, th.T.TextSize.Scale(((42.0/16.0)/100)*th.TextSize), txt)
}

func H2(th *Theme, txt string) material.LabelStyle {
	return material.Label(th.T, th.T.TextSize.Scale(((38.0/16.0)/100)*th.TextSize), txt)
}

func H3(th *Theme, txt string) material.LabelStyle {
	return material.Label(th.T, th.T.TextSize.Scale(((32.0/16.0)/100)*th.TextSize), txt)
}

func H4(th *Theme, txt string) material.LabelStyle {
	return material.Label(th.T, th.T.TextSize.Scale(((26.0/16.0)/100)*th.TextSize), txt)
}

func H5(th *Theme, txt string) material.LabelStyle {
	return material.Label(th.T, th.T.TextSize.Scale(((22.0/16.0)/100)*th.TextSize), txt)
}

func H6(th *Theme, txt string) material.LabelStyle {
	return material.Label(th.T, th.T.TextSize.Scale(((18.0/16.0)/100)*th.TextSize), txt)
}

func Body(th *Theme, txt string) material.LabelStyle {
	return material.Label(th.T, th.T.TextSize.Scale(((16.0/16.0)/100)*th.TextSize), txt)
}

func Small(th *Theme, txt string) material.LabelStyle {
	return material.Label(th.T, th.T.TextSize.Scale(((12.0/16.0)/100)*th.TextSize), txt)
}
