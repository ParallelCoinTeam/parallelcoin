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
	methods:{}
			
			`,
		Template: `
<template>
<div id="container" v-show="!duoSystem.bios.isBoot" class="swrap display lightTheme">
<duoLogo></duoLogo>
<duoHeader></duoHeader>
<duoSidebar></duoSidebar>
<duoMain class="flx flc grayGrad duoMain">
<keep-alive><duoScreenX class="flx flc fii duoScreenX"></duoScreenX></keep-alive>
</duoMain>
</div></template>`,
		Css: `
.display{
display: grid;
grid-gap: 0;
grid-template-columns: 60px 1fr;
grid-template-rows: 60px 1fr;
grid-template-areas:
"Logo Header"
"Sidebar Main"
}

.duoMain{
	padding:15px;
}


.dashboardParent{
	display:block;
	width:100%;
	height:100%;
}
`}
}
