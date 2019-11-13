<style>
.e-card .e-card-header .e-card-header-image.balance {
	position:absolute;
	bottom:0;
	left:0;
}

.e-card.duoCard {
width: 100%;
height:100%;
background: linear-gradient(-135deg, #f4f4f4 0%, #d4d4d4) 100%;
}
.e-card.duoCard div.e-card-sub-title{
text-align:right;

}

.e-card.duoCard  .e-card-sub-title span{
	font-size:36px;
	font-weight:100;
	letter-spacing:-1px;
}
.e-card.duoCard .e-card-content{
margin-left:auto;
}

</style>

<div class="rwrap flx">
	<div class="e-card flx flc justifyBetween duoCard">
		<div class="e-card-header">
			<div class="e-card-header-caption">
				<div class="e-card-header-title">Balance:</div>
				<div class="e-card-sub-title"><span v-html="this.duOSys.status.balance.balance"></span> DUO</div>
			</div>
			<div class="e-card-header-image balance"></div>
		</div>
		<div class="flx flc e-card-content">
			<small><span>Pending: </span><strong><span v-html="this.duOSys.status.balance.unconfirmed"></span></strong></small>
			<small><span>Transactions: </span><strong><span v-html="this.duOSys.status.txsnumber"></span></strong></small>
		</div>
	</div>
</div>