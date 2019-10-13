package pnl

import "github.com/p9c/pod/cmd/gui/mod"

func ChartA() mod.DuOScomp {
	return mod.DuOScomp{
		IsApp:    true,
		Name:     "Chart A",
		ID:       "panelcharta",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "chart",
		Js: `
	data () { return { 
	duOSys }},
		`,
		Template: `<div class="rwrap">
Chart A
		</div>`,
		Css: `

		`,
	}
}

func ChartB() mod.DuOScomp {
	return mod.DuOScomp{
		IsApp:    true,
		Name:     "Chart B",
		ID:       "panelchartb",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "chart",
		Js: `
	data () { return { 
	duOSys }},
		`,
		Template: `<div class="rwrap">
Chart B
		</div>`,
		Css: `

		`,
	}
}
