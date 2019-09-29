package html

func AMPHTML(duos, x, lib, css, sw string) string {
	return `<!doctype html>
<html ⚡ lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width,minimum-scale=1,initial-scale=1">
    <link rel="canonical" href="//parallelcoin.io"/>
    <title>ParallelCoin - DUO - True Story"</title>
    <meta name="description" content="ParallelCoin">
<link rel="preconnect dns-prefetch" href="https://fonts.gstatic.com/" crossorigin>
` + css + `
<meta name="amp-experiments-opt-in" content="visibility-v2,visibility-v3">
<link href="https://fonts.googleapis.com/css?family=Open+Sans:100,300,400" rel="stylesheet" />
<link href="https://fonts.googleapis.com/css?family=IBM+Plex+Mono:400" rel="stylesheet" />
` + duos + `
  	</head>
  	<body>
	<header id="boot"></header>
` + x + `
</html>`
}
