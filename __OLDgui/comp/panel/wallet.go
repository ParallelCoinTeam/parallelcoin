package panel

import "github.com/p9c/pod/cmd/gui/OLDbkp/vue/mod"

func Wallet() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Wallet",
		ID:       "panelwallet",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "test",
		Js: `
	data (){
			return {
			phrase:         '',
			addr:           '',
			amount:         0,
			showWalletPass: false,
		}
	},
	watch: {
		showWalletPass: function(){
		Swal.fire({
		title: 'Enter your password',
		input: 'password',
		inputPlaceholder: 'Enter your password',
		inputAttributes: {
		maxlength: 10,
		autocapitalize: 'off',
		autocorrect: 'off'
		 }
	});

		if (this.pharse) {
		Swal.fire('Entered password: ' + this.pharse);
		this.fnsend();
		this.showWalletPass = false;
	}
	},
	},
	methods: {
		fnsend:               function(){ wallet.fnSend(this.phrase,this.addr, this.amount); },
		toggleShowWalletPass: function(){ this.showWalletPass = !this.showWalletPass; },
		toggleShowNewAddr:    function(){ this.showNewAddr = !this.showNewAddr;},
		goGetAddresses:       function(){ wallet.getAddresses(); },
		showNewAddr:          function(){ Swal.fire({ position: 'center', type: 'success', title: 'You created new address', text: shared.node.receive.addr, showConfirmButton: false, timer: 9999 })},
	},
`,
		Template: `<div class="rwrap">
<ul><li>
	<input v-model="addr" class="fii mrl pds"><input v-model.number="amount" type="number" class="pds">
<ejs-button type="submit" @click="showWalletPass = true" iconCss="e-btn-sb-icon e-send-icon">Send</ejs-button></li>
	</li><li>
	<ejs-button type="submit" @click="goGetAddresses(); showNewAddr();" iconCss="e-btn-sb-icon e-receive-icon">Submit</ejs-button></li>

	  </ul>
</div>`,
		Css: `
		.boot{
			background:red;
		}
		`,
	}
}
