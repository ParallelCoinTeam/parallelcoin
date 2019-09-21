package serv

import "github.com/parallelcointeam/parallelcoin/cmd/gui/vue/mod"

func Services() mod.DuoVUEcomp {
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
		this.goGetTransactions();
	},
	watch: {
		'this.duoSystem.alert.time': function(newVal, oldVal) {
			setTimeout(() => {
				this.$refs.defaultRef.show({
					title: 'Adaptive Tiles Meeting', content: 'Conference Room 01 / Building 135 10:00 AM-10:30 AM',
					icon: 'e-meeting'
					});
				},200);
			}
		},
	mounted: function() {
		setTimeout(() => {
			this.$refs.defaultRef.show({
				title: 'Adaptive Tiles Meeting', content: 'Conference Room 01 / Building 135 10:00 AM-10:30 AM',
				icon: 'e-meeting'
            	});
			},200);
	},
 	methods: { 
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
Template: `<div><ejs-toast ref='defaultRef' id='toast_default' :created="created" :position='position'></ejs-toast></div>`,
	}
}
