package sys

import "github.com/p9c/pod/cmd/gui/vue/mod"

func Dev() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Dev",
		ID:       "dev",
		Version:  "0.0.1",
		CompType: "core",
		SubType:  "dev",
		Js: `
		data () { return { 
		duoSystem,
			}}`,
		Template: `<dev class="swrap dev" v-show="this.duoSystem.bios.isDev">
		<div class="dev"><h1>Layout</h1>
		DADADAD
		</div>
		</dev>`,
		Css: `
		.dev{
			background:blue;
		}
		`,
	}
}
