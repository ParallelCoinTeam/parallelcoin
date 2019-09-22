package comp

import "github.com/parallelcointeam/parallelcoin/cmd/gui/vue/mod"

func Header() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    false,
		Name:     "Header",
		ID:       "duoHeader",
		Version:  "0.0.0",
		CompType: "comp",
		SubType:  "header",
		Js: `
	  data() {
    return {
		duoSystem,
	}},

`,
		Template: `
<template>
	<div id="head" class="posRel flx justifyBetween bgLight lineBottom duoHeader">
		<div class="searchContent">
			<div class="analysis">ParallelCoin</div>
		</div>
		<div class="right-content">
			<ejs-button
			ref="toggleFullScreenBtn"
			iconCss="e-btn-sb-icons e-play-icon"
			cssClass="e-small e-flat e-success"
			:isPrimary="true"
			:isToggle="true"
			v-on:click.native="btnClick"
			>FullScreen</ejs-button>

    <ejs-button id="toggle" ref="toggleBoardbtn" class="e-btn e-info"  cssClass="e-flat" iconCss="e-icons burg-icon" isToggle="true" v-on:click.native="btnClick">Open</ejs-button>

			</div>
		</div>
	</div>
</template>
		`,
		Css: `

`}
}
