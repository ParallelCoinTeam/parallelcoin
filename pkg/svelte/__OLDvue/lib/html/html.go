package html

import (
	lib2 "github.com/p9c/pod/pkg/svelte/__OLDvue/lib"
	css2 "github.com/p9c/pod/pkg/svelte/__OLDvue/lib/css"
)

var HTML = VUEHTML(VUEx(lib2.VUElogo(), VUEheader(), VUEnav(lib2.ICO()), ScreenOverview()), css2.CSS(css2.ROOT(), css2.GRID(), css2.COLORS(), css2.HELPERS(), css2.NAV()))

func VUEHTML(x, s string) string {
	return `
<!DOCTYPE html><html lang="en" >
	<head>
		<meta charset="UTF-8">
		<title>ParallelCoin Wallet - True Story</title>
		<style type="text/css">` + s + `</style>
	</head>
  	<body>
		<header id="boot"></header>
` + x + `
		<footer id="dev"></footer>
	</body>
</html>`
}
