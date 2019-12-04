Vue.component('PanelStatus', {
		name: 'PanelStatus',
		data () { return { 
			rc: rcvar }},
		template: `<div class="rwrap e-card flx flc justifyBetween duoCard">
		<ul class="rf flx flc noPadding justifyEvenly">
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Block Height: </span><strong class="rcx6"><span v-html="rc.osHeight"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Local Hashrate: </span><strong class="rcx6"><span v-html="rc.osHashes"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Network Hashrate: </span><strong class="rcx6"><span v-html="rc.osNetHash"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Difficulty: </span><strong class="rcx6"><span v-html="rc.osDifficulty"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Blocks: </span><strong class="rcx6"><span v-html="rc.osBlockCount"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Connections: </span><strong class="rcx6"><span v-html="rc.osConnections"></span></strong></li>
		</ul>
	</div>`,
});