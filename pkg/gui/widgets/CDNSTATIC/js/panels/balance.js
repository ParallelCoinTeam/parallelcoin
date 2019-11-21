const PanelBalance = {
		name: 'PanelBalance',
		props:{
			balance:Object,
			txsnumber:Object,
		},
		template: `<div id="panelwalletstatus" class="rwrap flx Balance">
			<div class="e-card flx flc justifyBetween duoCard">
				<div class="e-card-header">
					<div class="e-card-header-caption">
						<div class="e-card-header-title">Balance:</div>
						<div class="e-card-sub-title"><span v-html="this.balance.balance"></span> DUO</div>
					</div>
					<div class="e-card-header-image balance"></div>
				</div>
				<div class="flx flc e-card-content">
					<small><span>Pending: </span><strong><span v-html="this.balance.unconfirmed"></span></strong></small>
				<small><span>Transactions: </span><strong><span v-html="this.txsnumber"></span></strong></small>
				</div>
			</div>
		</div>`,
}