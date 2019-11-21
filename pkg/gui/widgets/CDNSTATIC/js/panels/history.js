const PanelHistory = {
	name: 'PanelHistory',
	data () { return { 
		pageSettings: { pageSize: 10, pageSizes: [10,20,50,100], pageCount: 5 },
		ddldata: ['All', 'generated', 'sent', 'received', 'immature']
	}},
	methods: {
	},
	props:{
		transactions:Object,
	},
	template: `<div class="rwrap">
		<div class="select-wrap">
			<ejs-dropdownlist id='ddlelement' :dataSource='ddldata' placeholder='Select category to filter'></ejs-dropdownlist>
		</div>
		<ejs-grid :dataSource="this.transactions.txsEx.txs" height="100%" :allowPaging="true" :pageSettings='pageSettings'>
			<e-columns>
				<e-column field='category' headerText='Category' textAlign='Right' width=90></e-column>
				<e-column field='time' headerText='Time' format='unix'  textAlign='Right' width=90></e-column>
				<e-column field='txid' headerText='TxID' textAlign='Right' width=240></e-column>
				<e-column field='amount' headerText='Amount' textAlign='Right' width=90></e-column>
			</e-columns>
		</ejs-grid>
	</div>`,
}