package html

func VUEnav(ico map[string]string) string {
	return `<ul id="menu" class="lsn noPadding">
	<li id='menuoverview' class='sidebar-item current'>
		<button onclick="external.invoke('overview')" class="noMargin noPadding noBorder bgTrans sXs">` + ico["overview"] + `</button>
	</li>
	<li id='menutransactions' class='sidebar-item'>
		<button onclick="external.invoke('hisotry')" class="noMargin noPadding noBorder bgTrans sXs">` + ico["history"] + `</button>
	</li>
	<li id='menuaddressbook' class='sidebar-item'>
		<button onclick="external.invoke('addressbook')" class="noMargin noPadding noBorder bgTrans sXs">` + ico["addressbook"] + `</button>
	</li>
	<li id='menublockexplorer' class='sidebar-item'>
		<button onclick="external.invoke('overview')" class="noMargin noPadding noBorder bgTrans sXs">` + ico["overview"] + `</button>
	</li>
	<li id='menusettings' class='sidebar-item'>
		<button onclick="external.invoke('settings')" class="noMargin noPadding noBorder bgTrans sXs">` + ico["settings"] + `</button>
	</li>
</ul>`
}
