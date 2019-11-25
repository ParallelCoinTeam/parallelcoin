Vue.component('PageOverview', {
	el: "#overview",
  	name: "Overview",
  	template: `<main class="pageOverview">
	  <div id="panelwalletstatus" class="Balance">
	  <PanelBalance />
	  </div>
	  <div id="panelsend" class="Send">
	  <PanelSend />
	  </div>
	  <div id="panelnetworkhashrate" class="NetHash">
	  <PanelNetworkHashrate />
	  </div>
	  <div id="panellocalhashrate" class="LocalHash">
	  <PanelLocalHashrate />
	  </div>
	  <div id="panelstatus" class="Status">
		<PanelStatus />
	  </div>
	  <div id="paneltxsex" class="Txs">
	  <PanelLatestTx />
	  </div>
	  <div class="Log">
	  7
	  </div>
	  <div class="Info">
	  8
	  </div>
	  <div class="Time">
	  9
	  </div>
  </main>`,
});