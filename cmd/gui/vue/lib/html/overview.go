package html

import "github.com/p9c/pod/cmd/gui/mod"

func screenOverview() mod.DuOScomp {
	return mod.DuOScomp{
		IsApp:    false,
		Name:     "Screen Overview",
		ID:       "screenOverview",
		Version:  "0.0.1",
		CompType: "screen",
		SubType:  "screen",
		Js: `
	data () { return {
		duOSys,
}},
 methods: { 
  } 
`,
		Template: `<div id="screenoverview" class="Overview"><div id="panelwalletstatus" class="Balance"></div><div id="panelsend" class="Send"></div><div id="panelnetworkhashrate" class="NetHash"></div><div id="panellocalhashrate" class="LocalHash"></div><div id="panelstatus" class="Status"></div><div id="paneltxsex" class="Txs"></div><div class="Log"></div><div class="Info"></div><div class="Time"></div></div>`,
	}
}
