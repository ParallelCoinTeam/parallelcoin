<script>
	import Logo from "%/components/ico/Logo.svelte";
	import Button from "./layout/Button.svelte";
	import Header from "./layout/Header.svelte";
	import Nav from "./layout/Nav.svelte";
	import Overview from "./layout/Overview.svelte";
</script>

<style>
  main {
    font-family: sans-serif;
    text-align: center;
  }
</style>

	<x id="x" v-show="!this.duOSys.bios.isBoot" class="bgDark lightTheme"><div id="display">
    <div class="grid-container bgDark">
    	<div class="flx fii Logo"><Logo /></div>
    	<Header />
    	<div class="Sidebar bgLight">
    		<div class="Open"></div>
    		<Nav />
    		<div class="Side"></div>
    	</div>
    	<div id="main" class="grayGrad Main"><Overview /></div>
    </div>
    </div>
</x>



