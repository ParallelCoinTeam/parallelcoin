<script>
    import { quintOut } from 'svelte/easing';
	import { fade, draw, fly } from 'svelte/transition';
	import { expand } from './transitions.js';
	import { inner, outer } from './shape.js';
  

  let progress = 0;
	let visible = true;

  function next() {
    setTimeout(() => {
      if (progress === 100) {
        progress = 0;
      }
      progress += 1;
      next();
    }, 100);
  }
  next();


//   const { preloading, page } = stores();

</script>

<label>
	<input type="checkbox" bind:checked={visible}>
	toggle me
</label>

{#if visible}

<div class="flx centered baseMargin plan eightWidth justifyBetween" out:fly="{{y: -20, duration: 999}}">
		{#each 'PLAN 9 FROM CRYPTO SPACE' as char, i}
			<span class="font9"
					in:fade="{{delay: 999 + i * 150, duration: 999}}"
			>{char}</span>
		{/each}
	</div>

	<svg id ="bootlogo" class="center" viewBox="0 0 108 128">
		<g out:fade="{{ duration: 999}}" opacity=0.2>
			<path
				in:draw="{{delay: 999 , duration: 9999}}"
				style="stroke:#cfcfcf; stroke-width: 1.5"
				d={inner}
			/>
		</g>
	</svg>
	<div class="centered name" out:fly="{{y: -20, duration: 999}}">
		{#each 'ParallelCoin' as char, i}
			<span
				in:fade="{{delay: 999 + i * 150, duration: 999}}"
			>{char}</span>
		{/each}
	</div>

{/if}


<div class="progress justifyCenter textCenter txDark">
	<caption>{progress}%</caption>
</div>

	

<style>
	#bootlogo {
		height: 38vh;
		width: auto;
	}

	path {
		fill: #303030;
		opacity: 1;
	}

	label {
		position: absolute;
		top: 1em;
		left: 1em;
	}

	.centered {
		letter-spacing: 0.12em;
		color: #cfcfcf;
		font-weight: 100;
	}
	.plan {
		font-size: 3vw;
	}
	.name {
		font-size: 9vw;
	}
	.centered span {
		will-change: filter;
	}
	.progress {
		position: absolute;
		bottom: 0;
	}

</style>