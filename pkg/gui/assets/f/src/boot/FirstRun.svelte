<script>
	import { fade, fly } from 'svelte/transition';
  import Toggle from "../com/Toggle.svelte";
  import { isFirstRun} from '../firmware.js';

	let isShowModal = false;

  let privPassphrase = "";
  let confirmPrivPassphrase = "";
  let pubPassphrase = "";
  let seed = "";
  let walletDir = "";
  let error = {};

 let autoPubSeedFile = true;

  
  const handleSubmit = () => {
    if (privPassphrase.trim() === "") {
      error.privPassphrase = "Passphrase  field is required";
      return;
    }
    if (confirmPrivPassphrase.trim() === "") {
      error.confirmPrivPassphrase = "Confirm Passphrase field is required";
      return;
    }
    if (confirmPrivPassphrase.trim() !== privPassphrase.trim()) {
      error.privPassphrase = "Passphrase fields do not match";
      error.confirmPrivPassphrase = "Passphrase fields do not match";
      return;
    }
    duos.createWallet(privPassphrase, seed, pubPassphrase, walletDir);
    isFirstRun.set(false);
    isShowModal = false;
  }


  // function endFirstRun(){
  //   duos.createWallet(privPassphrase, seed, pubPassphrase, walletDir);
  //   isFirstRun.set(false);
  // }


   function autoSBF() {
		autoPubSeedFile = !autoPubSeedFile;
    }

   function showModal() {
		isShowModal = !isShowModal;
    }
	

</script>

<div class="flx flc itemsCenter form">

<svg xmlns="http://www.w3.org/2000/svg" id="parallel" viewBox="0 0 108 128" width="108" height="128"><path style="fill:#cfcfcf;fill-opacity:1" d="M77.08,2.55c3.87,1.03 6.96,2.58 10.32,4.64c5.93,3.87 10.58,8.51 14.19,14.71c3.87,6.19 5.42,13.16 5.42,20.64c0,7.22 -1.81,14.18 -5.41,20.37c-3.61,6.45 -8.25,11.35 -14.19,14.96c-3.35,2.06 -6.96,3.87 -10.32,4.9c-3.87,1.03 -7.74,1.55 -11.61,1.55v-14.45c6.96,-0.26 13.42,-2.58 19.09,-8c5.67,-5.42 8.51,-11.87 8.51,-19.61c0,-7.74 -2.58,-14.19 -7.99,-19.6c-5.42,-5.42 -11.86,-8 -19.6,-8c-7.74,0 -14.44,2.58 -19.6,8c-5.42,5.42 -8,11.87 -8,19.6l0,85.9c-3.1,-3.1 -7.99,-7.74 -13.93,-13.67v-72.23c0,-3.87 0.52,-7.73 1.55,-11.35c1.03,-3.87 2.58,-7.22 4.64,-10.32c3.87,-5.93 8.52,-10.58 14.71,-14.45c6.19,-3.61 13.16,-5.16 20.64,-5.16c3.87,0 8,0.52 11.61,1.55zM78.37,42.28c0,7.22 -5.93,13.16 -13.15,13.16c-7.48,0.26 -13.16,-5.68 -13.16,-13.16c0,-7.22 5.94,-13.16 13.16,-13.16c7.22,0 13.15,5.93 13.15,13.16zM13.63,37.12l0,69.39c-6.19,-6.19 -11.09,-10.83 -13.93,-13.93l0,-55.46z" /></svg>
<h1 class="upCase txLight textCenter">Creating new wallet</h1>

  <span class="txLight textCenter baseMargin">Enter the private passphrase for your new wallet</span>
      <input
      class="fullWidth"
        type="password"
        placeholder="Passphrase"
        bind:value={privPassphrase}
        autocomplete="new-password" />
      {#if error.privPassphrase}
        <code>{error.privPassphrase}</code>
      {/if}
      <input
      class="fullWidth"
        type="password"
        placeholder="Confirm Passphrase"
        bind:value={confirmPrivPassphrase}
        autocomplete="new-password" />
      {#if error.confirmPrivPassphrase}
        <code>{error.confirmPrivPassphrase}</code>
      {/if}


    

        {#if !autoPubSeedFile}
        <div class="flx flc glsm" transition:fade>
          <span class="txLight textCenter marginTopSm">Public Passphrase</span>
          <input type="text" placeholder="pubPassphrase" bind:value={pubPassphrase} />
          {#if error.pubPassphrase}
            <code>{error.pubPassphrase}</code>
          {/if}

          <span class="txLight textCenter marginTopSm">Seed</span>
          <input type="text" placeholder="Seed" bind:value={seed} />
          {#if error.seed}
            <code>{error.seed}</code>
          {/if}

          <span class="txLight textCenter marginTopSm">Wallet Directory</span>
          <input type="text" placeholder="walletDir" bind:value={walletDir} />
          {#if error.walletDir}
            <code>{error.walletDir}</code>
          {/if}
        </div>
      {/if}
      <button type="submit" on:click={showModal}>Submit</button>


</div>


    <Toggle checked={true}  on:change={autoSBF} />




{#if isShowModal}
<div class="fullScreen">
<div class="modal-background" on:click={showModal}></div>
<div in:fly="{{ y: 200, duration: 600 }}" out:fade class="modal" role="dialog" aria-modal="true">
		<h2>Write down...</h2>
    <h4>but, are you sure?</h4>
    <div class="fullWidth flx confirm">
    <button class="flx fii justifyCenter textCenter cancel" on:click={showModal}>cancel</button>
    <button class="flx fii justifyCenter textCenter create" on:click={handleSubmit}>Create wallet</button>
      </div>
  </div>
</div>
{/if}