<script>
	import { duoSystem } from '../firmware.js'
    import BootLogo from "./logo/BootLogo.svelte";
    import FirstRun from "./FirstRun.svelte";

  duoSystem.visible = true;

</script>

<div class="fullScreen flx flc justifyEvenly itemsCenter bgDark boot">
        <FirstRun />

    {#if duoSystem.bootLogo}
	<BootLogo />
    {/if}

</div>

