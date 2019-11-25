const PanelBalance = {
		name: 'PanelBalance',
		template: `<div id="panelwalletstatus" class="rwrap flx Balance">
			<div class="e-card flx flc justifyBetween duoCard">
				<div class="e-card-header">
					<div class="e-card-header-caption">
						<div class="e-card-header-title">Balance:</div>
						<div class="e-card-sub-title"><span v-html="rcvar.osBalance"></span> DUO</div>
					</div>
					<div class="e-card-header-image balance"></div>
				</div>
				<div class="flx flc e-card-content">
					<small><span>Pending: </span><strong><span v-html="rcvar.osUnconfirmed"></span></strong></small>
				<small><span>Transactions: </span><strong><span v-html="rcvar.osTxsNumber"></span></strong></small>
				</div>
			</div>
		</div>`,
}