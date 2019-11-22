const PanelSettings = {
	name: 'PanelSettings',
	data () { return { 
		duoSystem }},
		created: function() {
			
			},
		methods: { 
			 
		}, 
		template: `<div class="rwrap">
		<div v-html="this.duoSystem.config.daemon.schema"></div>
		 <vue-form-generator class="flx flc fii" :schema="this.duoSystem.config.daemon.schema" :model="this.duoSystem.config.daemon.config"></vue-form-generator>
				</div>`,
}