package sys

import "github.com/p9c/pod/cmd/gui/vue/mod"

func Boot() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "BOOT",
		ID:       "boot",
		Version:  "0.0.1",
		CompType: "core",
		SubType:  "boot",
		Js:       "",
		Template: `<div class="rwrap boot" v-show="duoSystem.isBoot"><h1>Boot</h1></div>`,
		Css: `
		.boot{
			background:red;
		}
		`,
	}
}
