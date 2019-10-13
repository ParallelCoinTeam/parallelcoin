package comp

import "github.com/p9c/pod/cmd/gui/OLDbkp/vue/mod"

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
		cellAspectRatio:100/69,
	}},
`,
		Template: `
<template>
 <div><div class="dashboardParent">
              <ejs-dashboardlayout 
			v-for="(screen, key) in system.data.conf.display.screens" 
			v-show="duoSystem.isScreen === key" 
			:ref="'DashbordInstance' + key" 
			:columns="screen.columns" 
			:id="'Layout' + key" 
			:allowResizing="screen.allowResizing"
			:allowDragging="screen.allowDragging"
			:allowFloating="screen.allowFloating"
			:cellAspectRatio="cellAspectRatio"
			:cellSpacing=[15,15]>
			<e-panels>
				<e-panel 
					v-for="panel in screen.panels" 
					:sizeX="panel.sizeX" 
					:sizeY="panel.sizeY" 
					:row="panel.row" 
					:col="panel.col" 
					:header="panel.header"
					:cssClass="panel.cssClass" 
					:content="panel.content">
				</e-panel>
            </e-panels>

    </ejs-dashboardlayout>            
          </div></div>
</template>
`,
		Css: `
`}
}
