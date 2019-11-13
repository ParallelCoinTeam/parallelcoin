<script>
   import { isPage } from '../../stores/pages.js';
    import { fade } from "svelte/transition";
	import Logo from "../ico/Logo.svelte";
    import Header from "./Header.svelte";
    import Nav from "./Nav.svelte";



let page;

	const unsubscribe = isPage.subscribe(value => {
		page = value;
	});






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

<div id="main" class="grayGrad Main" transition:fade={{ duration: 300 }}>

                <svelte:component this={$isPage}/>

</div>


    </div>
    </div>
</div>





<style>

</style>
