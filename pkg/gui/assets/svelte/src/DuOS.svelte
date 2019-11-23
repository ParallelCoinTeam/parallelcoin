<script>
    import {isFirstRun, isLogoVisible} from './firmware.js';
    import BootLogo from './boot/logo/BootLogo.svelte';
    import FirstRun from "./boot/FirstRun.svelte";


   let firstRun;
    const unsubscribeFirstRun = isFirstRun.subscribe(value => {
            firstRun = value;
    });
   
   
   let logoVisible;
    const unsubscribeLogoVisible = isLogoVisible.subscribe(value => {
            logoVisible = value;
    });
   

    




</script>


<div class="fullScreen flx flc justifyEvenly itemsCenter bgDark boot">
     



     {#if firstRun}
    <FirstRun autoPubSeedFile={true} />
    {/if}
    {#if !firstRun}
	    <BootLogo visible={logoVisible} />
    {/if}
</div>
