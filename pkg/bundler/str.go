package bnd

func Assets() DuOSassets {
	assets := make(map[string]DuOSasset)
	assets["indexhtml"] = DuOSasset{
		Name:        "index.html",
		ContentType: "text/html",
	}
	assets["sveltejs"] = DuOSasset{
		Sub:         "svelte",
		Name:        "svelte.js",
		ContentType: "application/javascript",
	}
	assets["sveltecss"] = DuOSasset{
		Sub:         "svelte",
		Name:        "svelte.css",
		ContentType: "text/css",
	}
	assets["cssroot"] = DuOSasset{
		Sub:         "css",
		Name:        "root.css",
		ContentType: "text/css",
	}
	assets["csscolors"] = DuOSasset{
		Sub:         "css",
		Name:        "colors.css",
		ContentType: "text/css",
	}
	assets["csshelpers"] = DuOSasset{
		Sub:         "css",
		Name:        "helpers.css",
		ContentType: "text/css",
	}
	assets["cssgrid"] = DuOSasset{
		Sub:         "css",
		Name:        "grid.css",
		ContentType: "text/css",
	}
	assets["bariolregular"] = DuOSasset{
		Sub:         "font",
		Name:        "bariolregular.ttf",
		ContentType: "application/x-font-ttf",
	}
	assets["bariolthin"] = DuOSasset{
		Sub:         "font",
		Name:        "bariolthin.ttf",
		ContentType: "application/x-font-ttf",
	}
	assets["bariolbold"] = DuOSasset{
		Sub:         "font",
		Name:        "bariolbold.ttf",
		ContentType: "application/x-font-ttf",
	}
	assets["bariolitalic"] = DuOSasset{
		Sub:         "font",
		Name:        "bariolitalic.ttf",
		ContentType: "application/x-font-ttf",
	}
	return assets
}
