import App from "./DuOS.svelte";
import Boot from "./boot/Boot.svelte";


const boot = new Boot({
    target: document.body,
    data: {
        entities: {}
    }
});

export default boot;



const app = new App({
    target: document.body,
    data: {
        entities: {}
    }
});

export default app;


