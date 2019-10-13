package panel

import "github.com/p9c/pod/cmd/gui/OLDbkp/vue/mod"

func ChartA() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Chart A",
		ID:       "panelcharta",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "chart",
		Js: `
	data () { return { 
	duoSystem }},
		`,
		Template: `<div class="rwrap">
Chart A
		</div>`,
		Css: `

		`,
	}
}

func ChartB() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Chart B",
		ID:       "panelchartb",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "chart",
		Js: `
	data () { return { 
	duoSystem }},
		`,
		Template: `<div class="rwrap">
Chart B
		</div>`,
		Css: `

		`,
	}
}
