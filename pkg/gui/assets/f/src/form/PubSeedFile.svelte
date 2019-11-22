<script>
  import { CurrentForm, PubSeedFileForm } from "../store.js";

  let pubPassphrase = "";
  let seed = "";
  let walletDir = "";
  let error = {};

  const handleSubmit = () => {
    if (pubPassphrase.trim() === "") {
      error.pubPassphrase = "Passphrase field must not be empty";
      return;
    }

    if (seed.trim() === "") {
      error.seed = "Seed field must not be empty";
      return;
    }
    if (walletDir.trim() === "") {
      error.walletDir = "walletDir field must not be empty";
      return;
    }
    Passphrase.setValues(pubPassphrase, seed, walletDir);
    CurrentForm.setValues();
  };
</script>


  <form class="form" on:submit|preventDefault={handleSubmit}>
    <div>
      <label>Public Passphrase</label>
      <input type="text" placeholder="pubPassphrase" bind:value={pubPassphrase} />
      {#if error.pubPassphrase}
        <code>{error.pubPassphrase}</code>
      {/if}
    </div>
    <div>
      <label>Seed</label>
      <input type="text" placeholder="Seed" bind:value={seed} />
      {#if error.seed}
        <code>{error.seed}</code>
      {/if}
    </div>
    
    <div>
      <label>Wallet Directory</label>
      <input type="text" placeholder="walletDir" bind:value={walletDir} />
      {#if error.walletDir}
        <code>{error.walletDir}</code>
      {/if}
    </div>
    <div>
      <button type="submit">Create wallet</button>
    </div>
  </form>
