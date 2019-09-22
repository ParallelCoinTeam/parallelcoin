package comp

import "github.com/parallelcointeam/parallelcoin/cmd/gui/vue/mod"

func Screen() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    false,
		Name:     "Screen",
		ID:       "duoScreenX",
		Version:  "0.0.1",
		CompType: "core",
		SubType:  "screen",
		Js: `
	  data() {
    return {
		duoSystem,
		layout: duoSystem.config.display.screens[duoSystem.activeLayout],
	}},
`,
		Template: `
<template>
 <div class="dashboardParent">   
            <ejs-dashboardlayout id='analysisLayout' :columns='layout.columns' ref='analysisLayout' :cellSpacing='layout.cellSpacing' :cellAspectRaito='layout.aspectRatio'>
                <e-panels>
                    <e-panel v-for="panel in layout.panels" :sizeX="panel.sizeX" :sizeY="panel.sizeY" :row="panel.row" :col="panel.col" :content="panel.content"></e-panel>
                </e-panels>
            </ejs-dashboardlayout>            
          </div>
</template>
`,
		Css: `
`}
}

