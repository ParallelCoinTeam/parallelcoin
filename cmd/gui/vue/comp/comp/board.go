package comp

import "github.com/p9c/pod/cmd/gui/vue/mod"

func Board() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    false,
		Name:     "Board",
		ID:       "duoBoard",
		Version:  "0.0.0",
		CompType: "comp",
		SubType:  "board",
		Js: `
	  data() {
    return {
		duoSystem,
	}},

`,
		Template: `
<template>
<div>
<ejs-sidebar id="default-sidebar" ref="sidebar" :type="type" :target="target" :position="position" :enablePersistence="enablePersistence">
      <div class="title"> Sidebar content</div>
       <div class="sub-title">
        Click the button to close the Sidebar.
    </div>
    <div class="center-align">
    </div>
</ejs-sidebar></div>
</template>
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
