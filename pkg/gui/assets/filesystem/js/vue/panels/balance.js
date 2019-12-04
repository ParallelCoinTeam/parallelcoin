Vue.component('PanelBalance', {
		name: 'PanelBalance',
		data () { return { 
			rc: rcvar }},
		template: `<div id="panelwalletstatus" class="rwrap flx">
			<div class="e-card flx flc justifyBetween duoCard">
				<div class="e-card-header">
					<div class="e-card-header-caption">
						<div class="e-card-header-title">Balance:</div>
						<div class="e-card-sub-title"><span v-html="rc.osBalance"></span> DUO</div>
					</div>
					<div class="e-card-header-image balance"></div>
				</div>
				<div class="flx flc e-card-content">
					<small><span>Pending: </span><strong><span v-html="rc.osUnconfirmed"></span></strong></small>
				<small><span>Transactions: </span><strong><span v-html="rc.osTxsNumber"></span></strong></small>
				</div>
			</div>
		</div>`,
});