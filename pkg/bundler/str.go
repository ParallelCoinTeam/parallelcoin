package bnd

func Assets() DuOSassets {
	assets := make(map[string]DuOSasset)
	assets["index.html"] = DuOSasset{
		Sub:         "html",
		Name:        "index.html",
		ContentType: "text/html",
	}
	assets["svelte.js"] = DuOSasset{
		Sub:         "sve",
		Name:        "svelte.js",
		ContentType: "application/javascript",
	}
	assets["svelte.css"] = DuOSasset{
		Sub:         "sve",
		Name:        "svelte.css",
		ContentType: "text/css",
	}
	assets["root.css"] = DuOSasset{
		Sub:         "css",
		Name:        "root.css",
		ContentType: "text/css",
	}
	assets["colors.css"] = DuOSasset{
		Sub:         "css",
		Name:        "colors.css",
		ContentType: "text/css",
	}
	assets["helpers.css"] = DuOSasset{
		Sub:         "css",
		Name:        "helpers.css",
		ContentType: "text/css",
	}
	assets["grid.css"] = DuOSasset{
		Sub:         "css",
		Name:        "grid.css",
		ContentType: "text/css",
	}
	assets["bariolregular.ttf"] = DuOSasset{
		Sub:         "fonts",
		Name:        "bariolregular.ttf",
		ContentType: "application/x-font-ttf",
	}
	return assets
}
