package sys

import "github.com/p9c/pod/pkg/duos/mod"

func Display() mod.DuOScomp {
	return mod.DuOScomp{
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
<x id="container" v-show="!this.duoSystem.bios.isBoot" class="x lightTheme">
<duoLogo></duoLogo>
<duoHeader></duoHeader>
<duoSidebar></duoSidebar>
<duoMain class="flx flc grayGrad duoMain">
<duoScreenX class="duoScreenX"></duoScreenX>
</duoMain>
</x></template>`,
		Css: `


.duoMain{
	padding:15px;
	overflow:hidden;
	overflow-y:auto;
}


.duoScreenX{
}


.dashboardParent{
	position:relative;
	display:block;
	width:100%;
}
`}
}
