const PanelPeers = {
	name: 'PanelPeers',
	data () { return { 
		pageSettings: { pageSize: 5 }
	}},
		template: `<div class="rwrap">
		<ejs-grid :dataSource="rcvar.peers" :allowPaging="true" :pageSettings='pageSettings'>
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
}