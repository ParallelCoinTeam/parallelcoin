import DuOS from "./DuOS.svelte";



const duos = new DuOS({
    target: document.body,
    data: {
        entities: {}
    }
});

export default duos;


