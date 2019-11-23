<script>
	export let checked = false;
</script>

<div class="toggle">
      <span class="txLight textCenter baseMargin">Advanced settings</span>
   <label class="small">
     <input class="toggle-input" type="checkbox" bind:checked on:change>
     <span class="switch-left">Off</span>
     <span class="switch-right">On</span>
   </label>
</div>
