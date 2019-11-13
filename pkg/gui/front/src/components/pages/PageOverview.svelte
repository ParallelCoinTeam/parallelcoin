<script>
    import PanelBalance from "../panels/PanelBalance.svelte";
    import PanelSend from "../panels/PanelSend.svelte";
    import PanelNetworkHashrate from "../panels/PanelNetworkHashrate.svelte";
    import PanelLocalHashrate from "../panels/PanelLocalHashrate.svelte";
    import PanelStatus from "../panels/PanelStatus.svelte";
    import PanelLatestTx from "../panels/PanelLatestTx.svelte";

    import { status } from "../../stores/status.js";
</script>
<style>

</style>



<svelte:head>
	<title>Overview</title>
</svelte:head>



        <main class="pageOverview">


        <div id="panelwalletstatus" class="Balance">
            <PanelBalance balance={status.balance}/>
        </div>
        <div id="panelsend" class="Send">
            <PanelSend  />
        </div>
        <div id="panelnetworkhashrate" class="NetHash">
            <PanelNetworkHashrate  hashrate={status.networkhashrate}/>
        </div>
        <div id="panellocalhashrate" class="LocalHash">
            <PanelLocalHashrate hashrate={status.networkhashrate} />
        </div>
        <div id="panelstatus" class="Status">
            <PanelStatus status={status} />
        </div>
        <div id="paneltxsex" class="Txs">
            <PanelLatestTx />
        </div>
        <div class="Info">
        </div>
        <div class="Time">
        </div>

        </main>



<!-- 
   <Card.Card id="panelwalletstatus" class="Balance">
            <PanelBalance :balance="duoSystem.balance" :status="duoSystem.status"/>
        </Card.Card>
        <Card.Card id="panelsend" class="Send">
            PanelSend :balance="duoSystem.balance"
        </Card.Card>
        <Card.Card id="panelnetworkhashrate" class="NetHash">
            PanelNetworkHashrate :balance="duoSystem.balance"
        </Card.Card>
        <Card.Card id="panellocalhashrate" class="LocalHash">
            PanelLocalHashrate :balance="duoSystem.balance"
        </Card.Card>
        <Card.Card id="panelstatus" class="Status">
            PanelStatus :status="duoSystem.status"/
        </Card.Card>
        <Card.Card id="paneltxsex" class="Txs">
            PanelLatestTx :transactions="duoSystem.transactions"
        </Card.Card>
        <Card.Card class="Log">
        </Card.Card>
        <Card.Card class="Info">
        </Card.Card>
        <Card.Card class="Time">
        </Card.Card> -->