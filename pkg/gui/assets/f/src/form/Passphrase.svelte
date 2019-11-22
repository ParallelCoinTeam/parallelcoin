<script>
  import { PassphraseForm, CurrentForm } from "../store.js";

  import Toggle from '../com/Toggle.svelte'

  let privPassphrase = "";
  let confirmPrivPassphrase = "";
  let autoPubSeedFile = true;
  let error = {};

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



    PassphraseForm.setValues(privPassphrase);
    CurrentForm.success();

  };

   function autoSBF() {
      autoPubSeedFile++
    }
</script>



  <form class="form flx flc" on:submit|preventDefault={handleSubmit}>
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
      <Toggle checked={true} bind:autoPubSeedFile on:change={autoSBF} />
      <button type="submit">Submit</button>
  </form>
