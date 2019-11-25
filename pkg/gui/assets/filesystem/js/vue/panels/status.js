Vue.component('PanelStatus', {
		name: 'PanelStatus',
		template: `<div class="rwrap">
		<ul class="rf flx flc noPadding justifyEvenly">
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Version: </span><strong class="rcx6"><span v-html="rcvar.osStatus.ver"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Uptime: </span><strong class="rcx6"><span v-html="rcvar.osStatus.uptime"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Memory: </span><strong class="rcx6"><span v-html="rcvar.osStatus.mem.total"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Disk: </span><strong class="rcx6"><span v-html="rcvar.osStatus.disk.total"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Chain: </span><strong class="rcx6"><span v-html="rcvar.osStatus.net"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Blocks: </span><strong class="rcx6"><span v-html="rcvar.osStatus.blockcount"></span></strong></li>
			<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Connections: </span><strong class="rcx6"><span v-html="rcvar.osStatus.connectioncount"></span></strong></li>
		</ul>
	</div>`,
});