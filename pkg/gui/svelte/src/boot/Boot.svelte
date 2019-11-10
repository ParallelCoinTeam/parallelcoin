<script>
	import { duoSystem } from '../firmware.js'
    import Kernel from "./Kernel.svelte";
    import BootLogo from "./logo/BootLogo.svelte";

  duoSystem.visible = true;

</script>

<div class="fullScreen flx flc justifyEvenly itemsCenter bgDark boot">


    {#if duoSystem.bootLogo}
	<BootLogo />
    {/if}

</div>

