package sys

import "github.com/parallelcointeam/parallelcoin/cmd/gui/vue/mod"

func Display() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Display",
		ID:       "display",
		Version:  "0.0.0",
		CompType: "core",
		SubType:  "display",
		Js: `
	  data() {
    return {
		duoSystem,
	}},
	components:{
		duoLogo,
		duoHeader,
		duoSidebar,
		duoScreenX,
		
	},
			
			
			`,
		Template: `
<template>
<x id="container" v-show="!duoSystem.bios.isBoot" class="x lightTheme">
<duoLogo></duoLogo>
<duoHeader></duoHeader>
<duoSidebar></duoSidebar>
<duoMain class="flx flc grayGrad duoMain">
<duoScreenX class="flx flc fii duoScreenX"></duoScreenX>
</duoMain>
</x></template>`,
		Css: `


.duoMain{
	overflow:hidden;
	overflow-y:auto;
}


.duoScreenX{
	padding:15px;
}


.dashboardParent{
	position:relative;
	display:block;
	width:100%;
}
`}
}
