package serv

import "github.com/parallelcointeam/parallelcoin/cmd/gui/vue/mod"

func SrvNode() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Node Services",
		ID:       "srvnode",
		Version:  "0.0.1",
		CompType: "core",
		SubType:  "node",
		Js: `
	data () { return { 
		duoSystem,
		txFrom:0,
		txCount:10,
		txCat:'default'
	}},
	created: function() {
		this.getStatus();
		this.duoSystem.timer = setInterval(this.getStatus, 1000);
		external.invoke('addressBook');
		this.goGetTransactions();
	},
	watch: {}, 
	methods: { 
		getStatus: function(){
			external.invoke('status');	
		},
		goGetTransactions: function(){
			const txsCmd = {
			from: this.txFrom,
			count: this.txCount,
			cat:this.txCat,
			};
			const txsCmdStr = JSON.stringify(txsCmd);
			external.invoke('transactions:'+txsCmdStr);
		},
		cancelAutoUpdate: function() { clearInterval(this.duoSystem.timer); },
	}, 
		beforeDestroy() {
		clearInterval(this.duoSystem.timer)
	} 
`,
	}
}
