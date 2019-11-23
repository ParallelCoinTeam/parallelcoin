const PanelStatus = {
		name: 'PanelStatus',
		props:{
			status:Object,
		},
		template: `<div class="rwrap">
		<ul class="rf flx flc noPadding justifyEvenly">
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Version: </span><strong class="rcx6"><span v-html="this.status.ver"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Wallet version: </span><strong class="rcx6"><span v-html="this.status.walletver.podjsonrpcapi.versionstring"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Uptime: </span><strong class="rcx6"><span v-html="this.status.uptime"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Memory: </span><strong class="rcx6"><span v-html="this.status.mem.total"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Disk: </span><strong class="rcx6"><span v-html="this.status.disk.total"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Chain: </span><strong class="rcx6"><span v-html="this.status.net"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Blocks: </span><strong class="rcx6"><span v-html="this.status.blockcount"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Connections: </span><strong class="rcx6"><span v-html="this.status.connectioncount"></span></strong></li>
		</ul>
	</div>`,
}