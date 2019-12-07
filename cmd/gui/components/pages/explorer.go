package pages

import "github.com/p9c/pod/cmd/gui"


func PageExplorer()gui.DuOScomP{
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


Vue.component('PageExplorer', {
	data () { return { 
		rcvar }},
		name: 'Explorer',
		components: {},
		template: `<main class="pageExplorer">
        pageExplorer
    	</main>`
});