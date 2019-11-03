<script>
    import { fade } from "svelte/transition";
	import Card from 'components/Card';
	import Logo from "../com/ico/Logo.svelte";
    import Header from "../com/layout/Header.svelte";
    import Nav from "../com/layout/Nav.svelte";
	import PanelBalance from "../com/pnl/PanelBalance.svelte";
	import PanelSend from "../com/pnl/PanelSend.svelte";
	import PanelLatestTx from "../com/pnl/PanelLatestTx.svelte";
</script>

<svelte:head>
	<title>DuOS</title>
</svelte:head>


<div id="x" class="fullScreen bgDark flx lightTheme"><div id="display" class="fii">
    <div class="grid-container rwrap bgDark">
    	<div class="flx fii Logo"><Logo /></div>
    	<Header />
    	<div class="Sidebar bgLight">
    		<div class="Open"></div>
    		<Nav />
    		<div class="Side"></div>
    	</div>

		<component :is="duoSystem.isScreen"></component>

    	<main id="main" class="grayGrad Main" transition:fade={{ duration: 300 }}>
  <Card.Card id="panelwalletstatus" class="Balance">
   PanelBalance
  </Card.Card>
</main>


    </div>
    </div>
</div>

