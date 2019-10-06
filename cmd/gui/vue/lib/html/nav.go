package html

func VUEnav(ico map[string]string) string {
	return `<ul id="menu" class="lsn padZ">
	<li id='menuoverview' class='sidebar-item current'>
		<button onclick="external.invoke('overview')" class="marZ padZ">` + ico["overview"] + `</button>
	</li>
	<li id='menutransactions' class='sidebar-item'>
		<button onclick="external.invoke('hisotry')" class="marZ padZ">` + ico["history"] + `</button>
	</li>
	<li id='menuaddressbook' class='sidebar-item'>
		<button onclick="external.invoke('addressbook')" class="marZ padZ">` + ico["addressbook"] + `</button>
	</li>
	<li id='menublockexplorer' class='sidebar-item'>
		<button onclick="external.invoke('overview')" class="marZ padZ">` + ico["overview"] + `</button>
	</li>
	<li id='menusettings' class='sidebar-item'>
		<button onclick="external.invoke('settings')" class="marZ padZ">` + ico["settings"] + `</button>
	</li>
</ul>`
}
