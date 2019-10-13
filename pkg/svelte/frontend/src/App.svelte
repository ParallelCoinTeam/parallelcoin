<script>
	import Button from "./Button.svelte";
	import Logo from "./ico/Logo.svelte";
	import Header from "./Header.svelte";
	import Nav from "./Nav.svelte";
	import Overview from "./Overview.svelte";
</script>

<style>
  main {
    font-family: sans-serif;
    text-align: center;
  }
</style>

<main>
	<div id="x" v-show="!this.duOSys.bios.isBoot" class="bgDark lightTheme"><div id="display">
    <div class="grid-container bgDark">
    	<div class="flx fii Logo"><Logo /></div>
    	<Header />
    	<div class="Sidebar bgLight">
    		<div class="Open"></div>
    		<div class="Nav"><Nav /></div>
    		<div class="Side"></div>
    	</div>
    	<div id="main" class="grayGrad Main"><Overview /></div>
    </div>
    </div></div>
</main>