package panels

import "github.com/p9c/pod/cmd/gui"


func PanelStatus()gui.DuOScomP{
	return gui.DuOScomP{
		Name:        "",
		ID:          "",
		Version:     "",
		Description: "",
		State:       "",
		Image:       "",
		URL:         "",
		CompType:    "",
		SubType:     "",
		Js:          "",
		Html:        "",
		Css:         "",
	}
}

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
				<div class="flx fii baseMargin">
					<ul class="rf rwrap flx flc noPadding noMargin justifyEvenly">
						<li class="flx fwd justifyBetween lineBottom"><span>Block Height: </span><strong><span v-html="rc.osHeight"></span></strong></li>
						<li class="flx fwd justifyBetween lineBottom"><span>Local Hashrate: </span><strong><span v-html="rc.osHashes"></span></strong></li>
						<li class="flx fwd justifyBetween lineBottom"><span>Network Hashrate: </span><strong><span v-html="rc.osNetHash"></span></strong></li>
						<li class="flx fwd justifyBetween lineBottom"><span>Difficulty: </span><strong><span v-html="rc.osDifficulty"></span></strong></li>
						<li class="flx fwd justifyBetween lineBottom"><span>Blocks: </span><strong><span v-html="rc.osBlockCount"></span></strong></li>
						<li class="flx fwd justifyBetween lineBottom"><span>Connections: </span><strong><span v-html="rc.osConnections"></span></strong></li>
					</ul>
				</div>
			</div>
	</div>`,
});