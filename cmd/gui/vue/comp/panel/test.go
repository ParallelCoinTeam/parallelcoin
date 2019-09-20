package panel

import "github.com/parallelcointeam/parallelcoin/cmd/gui/vue/mod"

func Test() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Test",
		ID:       "paneltest",
		Version:  "0.0.1",
		CompType: "core",
		SubType:  "test",
		Js: `
	data () { return { 
	duoSystem }},
		`,
		Template: `<div class="rwrap">
<ul>				
<li>
<span>dddddd</span>
<hr><span>dddddd</span>
<hr><span v-html="this.duoSystem.screen"></span>
<hr>
</li>
</ul>
		</div>`,
		Css: `

		`,
	}
}
