package html

func VUEbalance() string {
	return `
<div class="rwrap">
	<div tabindex="0" class="e-card flx flc justifyBetween duoCard">
		<div class="e-card-header">
			<div class="e-card-header-caption">
				<div class="e-card-header-title">Balance:</div>
				<div class="e-card-sub-title"><span id="balanceValue">0</span> DUO</div>
			</div>
			<div class="e-card-header-image balance"></div>
		</div>
		<div class="flx flc e-card-content">
			<small><span>Pending: </span><strong><span id="unconfirmedValue">0</span></strong></small>
			<small><span>Transactions: </span><strong><span id="txsnumberValue">0</span></strong></small>
		</div>
    </div>
</div>
`
}
