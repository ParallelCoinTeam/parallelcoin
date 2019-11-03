<script>
import { duoSystem } from '../DuOSyS.js'
import PanelBalance from '../panels/PanelBalance.vue'
import PanelSend from '../panels/PanelSend.vue'
import PanelLatestTx from '../panels/PanelLatestTx.vue'
// import PanelLocalHashrate from '../panels/PanelLocalHashrate.vue'
// import PanelNetworkHashrate from '../panels/PanelNetworkHashrate.vue'
import PanelStatus from '../panels/PanelStatus.vue'

export default {
  name: 'Overview',
  data () { return { 
	duoSystem }},
  components: {
    PanelBalance,
    PanelSend,
    PanelLatestTx,
    // PanelLocalHashrate,
    // PanelNetworkHashrate,
    PanelStatus
  }
}
</script>

<template>
    <main class="pageOverview">
        <div id="panelwalletstatus" class="Balance">
            <PanelBalance :balance="duoSystem.balance" :status="duoSystem.status"/>
        </div>
        <div id="panelsend" class="Send">
            <PanelSend :balance="duoSystem.balance"/>
        </div>
        <div id="panelnetworkhashrate" class="NetHash">
            <!-- <PanelNetworkHashrate :balance="duoSystem.balance"/> -->
        </div>
        <div id="panellocalhashrate" class="LocalHash">
            <!-- <PanelLocalHashrate :balance="duoSystem.balance"/> -->
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
    </main>
</template>

<style>
@import '../../node_modules/@syncfusion/ej2-base/styles/material.css';
@import '../../node_modules/@syncfusion/ej2-buttons/styles/material.css';
</style>
