<script>
	import { duoSystem } from './firmware.js'
	import BootLogo from './boot/logo/BootLogo.svelte';
    import FirstRun from "./boot/FirstRun.svelte";


</script>


<div class="fullScreen flx flc justifyEvenly itemsCenter bgDark boot">
     {#if !duoSystem.bootLogo}
    <FirstRun />
    {/if}
    {#if duoSystem.bootLogo}
	<BootLogo></BootLogo>
    {/if}
</div>
