<script>
  import { duoSystem } from '../firmware.js'

  duoSystem.bootLogo = false;


let counter = 0;

const biosMenu = new BiosMenu({
  target: document.getElementById("biosMenu"),
  data: {
    count: duoSystem.counter
  }
});

const handler = setInterval(() => {
  counter += 1;

  biosMenu.set({ count: duoSystem.counter });

  if (counter === 9) {
    clearInterval(handler);

    setTimeout(() => {
    biosMenu.destroy();
    duoSystem.bootLogo = true;
    }, 1000);
  }
}, 1000);



</script>

<div class="fullScreen flx flc justifyCenter itemsCenter bgDark boot">
    <div class="biosMenu">
        <h1>{duoSystem.count}</h1>
        <h2>Press something for configuration</h2>
    </div>
</div>

