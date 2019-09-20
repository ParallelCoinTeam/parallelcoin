package comp

var GetAppHtml string = `<!DOCTYPE html><html lang="en" ><head><meta charset="UTF-8"><title>ParallelCoin Wallet - True Story</title></head><body>
<header is="boot" id="boot" v-show="duoSystem.isBoot"></header>
<section id="main" class="rwrap lightTheme">
<nav is="nav" id="nav"></nav>
<main :is="display" id="display"></main>
<aside id="alerts"></aside>
</section>
<section id="sub" class="rwrap hide">
<main :is="dev" id="dev"></main>
</section>
<footer is="serv" id="serv"></footer></body>`
