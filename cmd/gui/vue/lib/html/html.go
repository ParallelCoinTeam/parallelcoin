package html

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
