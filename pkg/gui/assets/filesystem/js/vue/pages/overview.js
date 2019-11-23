var PageOverview = {
	el: "#overview",
  	name: "Overview",
  	template: `<main class="pageOverview">
	  <div id="panelwalletstatus" class="Balance">
		  <PanelBalance :balance="duoSystem.status.balance" :txsnumber="duoSystem.status.txsnumber"/>
	  </div>
	  <div id="panelsend" class="Send">
		  <PanelSend />
	  </div>
	  <div id="panelnetworkhashrate" class="NetHash">
	<PanelNetworkHashrate :hashrate="duoSystem.status.networkhashrate"/>
	  </div>
	  <div id="panellocalhashrate" class="LocalHash">
	  <PanelLocalHashrate :hashrate="duoSystem.status.hashrate"/> 
	  </div>
	  <div id="panelstatus" class="Status">
		  <PanelStatus :status="duoSystem.status"/>
	  </div>
	  <div id="paneltxsex" class="Txs">
	  <PanelLatestTx :transactions="duoSystem.transactions"/>
	  </div>
	  <div class="Log">
	  </div>
	  <div class="Info">
	  </div>
	  <div class="Time">
	  </div>
  </main>`,
  	components: {
		PanelBalance,
		PanelSend,
		PanelLatestTx,
		PanelLocalHashrate,
		PanelNetworkHashrate,
		PanelStatus
	}
}