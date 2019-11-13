import {readable} from 'svelte/store';

export const bios = readable([], function start(set) {
    const interval = setInterval(() => {
        fetch(`http://127.0.0.1:3999/bios`)
            .then(resp => resp.json())
            .then(data => {
                set(data)
            });
    }, 100);
    return function stop() {
        clearInterval(interval);
    };
});
