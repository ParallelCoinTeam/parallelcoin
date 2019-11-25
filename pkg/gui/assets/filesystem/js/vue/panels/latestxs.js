Vue.component('PanelLatestTx', {
		name: 'PanelLatestTx',
		data () { return { 
			pageSettings: { pageSize: 10, pageCount: 1 }
		}},
	template: `<div class="rwrap">
			<ejs-grid :dataSource="rcvar.osLastTxs.txs" height="100%" :pageSettings='pageSettings'>
				<e-columns>
					<e-column field='category' headerText='Category' textAlign='Right' width=90></e-column>
					<e-column field='time' headerText='Time' format='unix'  textAlign='Right' width=90></e-column>
					<e-column field='txid' headerText='TxID' textAlign='Right' width=240></e-column>
					<e-column field='amount' headerText='Amount' textAlign='Right' width=90></e-column>
				</e-columns>
			</ejs-grid>
	</div>`,
});