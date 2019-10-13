import App from "./App.svelte";

const counter = document.querySelector('.counter');
// We use async/await because Go functions are asynchronous
const render = async () => {
  counter.innerText = `Count: ${await window.counterValue()}`;
};
btnIncr.addEventListener('click', async () => {
  await counterAdd(1); // Call Go function
  render();
});
btnDecr.addEventListener('click', async () => {
  await counterAdd(-1); // Call Go function
  render();
});
render();


const app = new App({
  target: document.body
});

export default app;
