import { writable } from "svelte/store";

function stepPassphrase() {
  const { subscribe, set, update } = writable({
    privPassphrase: "",
    autoPubSeedFile: false
  });

  const setValues = (privPassphrase, autoPubSeedFile) => {
    return update(state => ({
      ...state,
      privPassphrase,
      autoPubSeedFile
    }));
  };

  const logValues = () => {
    const unsubscribe = subscribe(value => {
      return console.log(value);
    });
  };

  return {
    subscribe,
    setValues,
    logValues
  };
}

function stepPubSeedFile() {
  const { subscribe, set, update } = writable({
    pubPassphrase: "",
    seed: "",
    walletDir: ""
  });

  const setValues = privPassphrase => {
    return update(state => ({
      ...state,
      pubPassphrase,
      seed,
      walletDir
    }));
  };

  const logValues = () => {
    const unsubscribe = subscribe(value => {
      return console.log(value);
    });
  };

  return {
    subscribe,
    setValues,
    logValues
  };
}

function currentForm() {
  const { subscribe, set, update } = writable({
    passPhrase: true,
    pubSF: false,
    success: false
  });

  const logValues = () => {
    const unsubscribe = subscribe(value => {
      return console.log(value);
    });
  };

  return {
    subscribe,
    logValues,
    setValues: () =>
      update(n => ({ ...n, passPhrase: !n.passPhrase, pubSF: !n.pubSF, success: false })),
    success: () =>
      update(n => ({ ...n, success: true, passPhrase: false, pubSF: false })),
      if (success){cw.createWallet(passPhrase, seed, pubPassphrase, walletDir)}
  };
}
export const PassphraseForm = stepPassphrase();
export const CurrentForm = currentForm();
export const PubSeedFileForm = stepPubSeedFile();
