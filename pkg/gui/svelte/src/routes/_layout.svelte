<script>
let boot =false;
  import './lang.js'
  import { stores } from "@sapper/app";
  import { onMount } from "svelte";
  import { fade } from "svelte/transition";
    import Boot from "../com/Boot.svelte";
	import UI from "../com/UI.svelte";
  import ProgressLinear from "components/ProgressLinear";
  export let segment;
  import {
    right,
    elevation,
    persistent,
    showNav,
    showNavMobile,
    breakpoint
  } from "smelte/src/stores.js";

  const { preloading, page } = stores();

</script>

<svelte:head>
  <title>ParallelCoin - The true story</title>
</svelte:head>


{#if boot}
	<Boot />
{/if}

{#if !boot}


{#if $preloading}
  <ProgressLinear app />
{/if}

<UI />
{/if}

<style>
@import "./../ass/css/theme.css";
</style>
