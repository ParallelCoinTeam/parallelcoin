<script>
	import { getClient } from '@matrx/svelte-realtime-store'
	import Boot from "./boot/Boot.svelte";

	const realtimeClient = getClient()
	const a = realtimeClient.realtime('a', 1000)
	const connected = realtimeClient.connected


</script>

<Boot />