<script>
	import Boot from './boot/Boot.svelte';
	import UI from './components/layout/UI.svelte';
   import { bios } from './stores/firmware.js';
</script>



{#if bios.boot}
<Boot />
{/if}
{#if !bios.boot}
<UI />
{/if}

