import {readable} from 'svelte/store';

export const lasTxs = readable([], function start(set) {
    const interval = setInterval(() => {
        fetch(`http://127.0.0.1:3999/lastxs`)
            .then(resp => resp.json())
            .then(data => {
                set(data)
            });
    }, 100);
    return function stop() {
        clearInterval(interval);
    };
});

