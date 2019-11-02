<script>
	import { locale, _ } from 'svelte-i18n'
    import { quintOut } from 'svelte/easing';
	import { fade, draw, fly } from 'svelte/transition';
	import { expand } from './com/intro/transitions.js';
	import { inner, outer } from './com/intro/shape.js';
  import ProgressLinear from 'components/ProgressLinear';
  import ProgressCircular from 'components/ProgressCircular';
  
  let progress = 0;
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
let visible = true;
</script>

<div class="fullScreen flx flc justifyBetween itemsCenter bgDark boot">
{#if visible}
	<svg id ="intrologo" class="marginTopBig" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 108 128">
		<g out:fade="{{duration: 900}}" opacity=0.2>
			<path
				in:draw="{{duration: 2000}}"
				style="stroke:#cfcfcf; stroke-width: 1.5"
				d={inner}
			/>
		</g>
	</svg>

	<div class="centered" out:fly="{{y: -20, duration: 800}}">
		{#each 'ParallelCoin' as char, i}
			<span
				in:fade="{{delay: 1000 + i * 150, duration: 800}}"
			>{char}</span>
		{/each}
	</div>
{/if}

<label>
	<input type="checkbox" bind:checked={visible}>
	toggle me
</label>


<div class="progress fullWidth textCenter txGray"><caption>{progress}%</caption>
<ProgressLinear {progress} class="bgPurple" /></div>

</div>

<style>
	#intrologo {
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
		font-size: 10vw;
		position: absolute;
		left: 50%;
		top: 72%;
		transform: translate(-50%,-50%);
		letter-spacing: 0.12em;
		color: #cfcfcf;
		font-weight: 100;
	}

	.centered span {
		will-change: filter;
	}

	.progress {
		position: absolute;
		bottom: 0;
	}

</style>