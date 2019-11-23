<script>
    import {isFirstRun} from './firmware.js';
    import BootLogo from './boot/logo/BootLogo.svelte';
    import FirstRun from "./boot/FirstRun.svelte";
		import { onMount } from 'svelte';


   let firstRun;

        const unsubscribe = isFirstRun.subscribe(value => {
            firstRun = value;
        });
        
    
    


</script>


<div class="fullScreen flx flc justifyEvenly itemsCenter bgDark boot">
     



     {#if firstRun}
    <FirstRun autoPubSeedFile={true} />
    {/if}
    {#if !firstRun}
	    <BootLogo />
    {/if}
</div>
