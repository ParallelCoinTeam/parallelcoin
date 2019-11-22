<script>
  import { PassphraseForm, CurrentForm, PubSeedFileForm } from "../store.js";

  import PassphraseComponent from "../form/Passphrase.svelte";
  import PubSeedFileComponent from "../form/PubSeedFile.svelte";

  console.log($PassphraseForm);
</script>

<style>
  .container {
    color: black;
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    transform: translateY(-80px);
  }
  .success {
    background-color: #9dff0087;
    border-radius: 10px;
    padding: 10px;
    border: 1px solid #00000033;
    color: green;
  }
  .success h1 {
    font-weight: bolder;
  }
</style>

<div class="flx flc">
<h1 class="upCase txLight textCenter">Creating new wallet</h1>
<span class="txLight textCenter baseMargin">Enter the private passphrase for your new wallet</span>
  {#if $CurrentForm.passPhrase}
    <PassphraseComponent />
  {:else if $CurrentForm.pubSF}
    <PubSeedFileComponent />
  {:else if $CurrentForm.success}
    <div class="success">
      <h1>Form submitted</h1>
    </div>
  {:else}
    <h3>Something is wrong here!</h3>
  {/if}

</div>
