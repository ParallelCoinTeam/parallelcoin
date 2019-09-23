package panel

import "github.com/p9c/pod/cmd/gui/vue/mod"

func WalletStatus() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Wallet status",
		ID:       "panelwalletstatus",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "status",
		Js: `
			data () { return { 
			duoSystem }}, 
		`,
		Template: `<div class="rwrap">

                           
    <ul class="rf flx flc noPadding justifyEvenly">
        <li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Balance: </span><strong class="rcx6"><span v-html="this.duoSystem.status.balance.balance"></span></strong></li>
        <li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Unconfirmed: </span><strong class="rcx6"><span v-html="this.duoSystem.status.balance.unconfirmed"></span></strong></li>
        <li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Transactions: </span><strong class="rcx6"><span v-html="this.duoSystem.status.txsnumber"></span></strong></li>
    </ul>
</div>`,
		Css: `



		`,
	}
}

func Status() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Status",
		ID:       "panelstatus",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "status",
		Js: `
			data () { return { 
			duoSystem }}, 
		`,
		Template: `<div class="rwrap">


    <ul class="rf flx flc noPadding justifyEvenly">
        <li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Version: </span><strong class="rcx6"><span v-html="this.duoSystem.status.ver"></span></strong></li>
		<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Wallet version: </span><strong class="rcx6"><span v-html="this.duoSystem.status.walletver.podjsonrpcapi.versionstring"></span></strong></li>

        <li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Uptime: </span><strong class="rcx6"><span v-html="this.duoSystem.status.uptime"></span></strong></li>


		<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Memory: </span><strong class="rcx6"><span v-html="this.duoSystem.status.mem.total"></span></strong></li>
		<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Disk: </span><strong class="rcx6"><span v-html="this.duoSystem.status.disk.total"></span></strong></li>

		<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Chain: </span><strong class="rcx6"><span v-html="this.duoSystem.status.net"></span></strong></li>
		<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Blocks: </span><strong class="rcx6"><span v-html="this.duoSystem.status.blockcount"></span></strong></li>
        <li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Connections: </span><strong class="rcx6"><span v-html="this.duoSystem.status.connectioncount"></span></strong></li>
    </ul>

</div>`,
		Css: `



		`,
	}
}

func NetworkHashRate() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Network Hashrate",
		ID:       "panelnetworkhashrate",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "status",
		Js: `
   data:function(){
return { 
			duoSystem,
  axes: [{
        startAngle: 270,
        endAngle: 90
    }],
    }},
		`,
		Template: `<div class='rwrap'><div class='wrapper'>
    <ejs-circulargauge width='300px' height='200px'>
           <e-axes>
            <e-axis>
              <e-pointers>
                <e-pointer :value='val'/>
              </e-pointers>
            </e-axis>
          </e-axes>
         </ejs-circulargauge>
</div></div>`,
		Css: `



		`,
	}
}
