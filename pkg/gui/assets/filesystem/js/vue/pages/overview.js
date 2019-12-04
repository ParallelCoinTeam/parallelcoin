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