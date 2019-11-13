<script>
    import BootLogo from "./logo/BootLogo.svelte";
    import { bios } from '../stores/firmware.js';

</script>
<div class="fullScreen flx flc justifyEvenly itemsCenter bgDark boot">
<h1 id="biosMessages">Boot</h1>

{#if bios.logo}
<BootLogo />
{/if}
</div>


