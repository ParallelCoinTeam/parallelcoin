package panel

import "github.com/parallelcointeam/parallelcoin/cmd/gui/vue/mod"

func Peers() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Peers",
		ID:       "panelpeers",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "status",
		Js: `
	data () { return { 
	duoSystem,
      pageSettings: { pageSize: 5 }
	}},
		`,
		Template: `<div class="rwrap">
        <ejs-grid :dataSource="this.duoSystem.peers" :allowPaging="true" :pageSettings='pageSettings'>
          <e-columns>
            <e-column field='addr' headerText='Address' textAlign='Right' width=90></e-column>
            <e-column field='pingtime' headerText='Ping time' width=120></e-column>
            <e-column field='bytessent' headerText='Sent' textAlign='Right' width=90></e-column>
			<e-column field='bytesrecv' headerText='Received' textAlign='Right' width=90></e-column>
			<e-column field='subver' headerText='Subversion' textAlign='Right' width=90></e-column>
			<e-column field='version' headerText='Version' textAlign='Right' width=90></e-column>
          </e-columns>
        </ejs-grid>
</div>`,
		Css: `
		
		`,
	}
}
