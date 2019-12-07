package pages

import "github.com/p9c/pod/cmd/gui"


func PageSettings()gui.DuOScomP{
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


Vue.use(VueFormGenerator);

Vue.component('PageSettings', {
	name: 'Settings',
		data () { return { 
		ds: duOSsettings }},
		template: `<main class="pageSettings"><div class="rwrap">
			<vue-form-generator class="flx flc fii" :schema="ds.data.daemon.schema" :model="ds.data.daemon.config"></vue-form-generator>
			<ejs-progressbutton content="Save Settings" :isPrimary="true" :spinSettings="spinRight"></ejs-progressbutton>
		</div></main>`
	});