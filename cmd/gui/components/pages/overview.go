package pages

import "github.com/p9c/pod/cmd/gui"


func PageLocalHashRate()gui.DuOScomP{
	return gui.DuOScomP{
		Name:        "",
		ID:          "",
		Version:     "",
		Description: "",
		State:       "",
		Image:       "",
		URL:         "",
		CompType:    "",
		SubType:     "",
		Js:          "",
		Html:        "",
		Css:         "",
	}
}


Vue.component('PageOverview', {
	el: "#overview",
  	name: "Overview",
	  template: `<main class="pageOverview">
	<div id="panelwalletstatus" class="Balance">
		<PanelBalance />
	</div>
	<div id="panelsend" class="SendReceive">
		<PanelSend />
	</div>
	<div id="paneltxsex" class="LastTx">
		<PanelLatestTx />
	</div>
	<div id="panelstatus" class="Status">
		<PanelStatus />
	</div>
  </main>`,
});