Vue.component('PageSettings', {
	name: 'Settings',
		data () { return { 
		duOSsettings }},
		template: `<main class="pageSettings"><div class="rwrap">
			<div v-html="duOSsettings.daemon.schema"></div>
			<vue-form-generator class="flx flc fii" :schema="duOSsettings.daemon.schema" :model="duOSsettings.daemon.config"></vue-form-generator>
		</div></main>`
	});