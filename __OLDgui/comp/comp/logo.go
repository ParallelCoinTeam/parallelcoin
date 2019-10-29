package comp

import "github.com/p9c/pod/cmd/gui/OLDbkp/vue/mod"

func Logo() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    false,
		Name:     "Logo",
		ID:       "duoLogo",
		Version:  "0.0.0",
		CompType: "comp",
		SubType:  "logo",
		Js: `
	  data() {
    return {
		duoSystem,
         type :'Push',
         target : '.content',
         position : 'Left',
         enablePersistence: true
	}},
`,
		Template: `
<template><div v-html="system.data.ico.logo" class="flx fii logo"></div></template>
		`,
		Css: `

.logo{
width:60px;
height:60px;
background-color:var(--dark);
}
.logofill {
fill:var(--light);
}
`}
}
