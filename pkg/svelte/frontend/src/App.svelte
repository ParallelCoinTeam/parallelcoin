<script>
	import { getClient } from '@matrx/svelte-realtime-store'

	const realtimeClient = getClient()
	const a = realtimeClient.realtime('a', 1000)
	const connected = realtimeClient.connected

	function handleA(event) {
		$a = $a + 1
	}
</script>

<h1>{$a}</h1>
<button on:click={handleA} class="button" disabled="{!$connected}">a++</button>