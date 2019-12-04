Vue.component('PanelStatus', {
		name: 'PanelStatus',
		data () { return { 
			rc: rcvar }},
		template: `<div class="rwrap e-card flx flc justifyBetween duoCard">
		<div class="e-card flx flc justifyBetween duoCard">
				<div class="e-card-header">
					<div class="e-card-header-caption">
						<div class="e-card-header-title">Status:</div>
					</div>
					<div class="e-card-header-image balance"></div>
				</div>
				<div class="cwrap flx e-card-content">
					<ul class="rf rwrap flx flc noPadding noMargin justifyEvenly">
						<li class="flx fwd spb htg rr justifyBetween"><span>Block Height: </span><strong><span v-html="rc.osHeight"></span></strong></li>
						<li class="flx fwd spb htg rr justifyBetween"><span>Local Hashrate: </span><strong><span v-html="rc.osHashes"></span></strong></li>
						<li class="flx fwd spb htg rr justifyBetween"><span>Network Hashrate: </span><strong><span v-html="rc.osNetHash"></span></strong></li>
						<li class="flx fwd spb htg rr justifyBetween"><span>Difficulty: </span><strong><span v-html="rc.osDifficulty"></span></strong></li>
						<li class="flx fwd spb htg rr justifyBetween"><span>Blocks: </span><strong><span v-html="rc.osBlockCount"></span></strong></li>
						<li class="flx fwd spb htg rr justifyBetween"><span>Connections: </span><strong><span v-html="rc.osConnections"></span></strong></li>
					</ul>
				</div>
			</div>
	</div>`,
});