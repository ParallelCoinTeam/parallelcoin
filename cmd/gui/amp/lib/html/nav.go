package html

func AMPnav() string {
	return `<ul id="menulist" class="e-listview">
	<li id='menuoverview' class='sidebar-item current' on-click='setScreen("overview","menuoverview");'>
		<span class='sf-icon-Dashboard list_svg'></span>
	</li>
	<li id='menutransactions' class='sidebar-item' on-click='setScreen("transactions","menutransactions");'>
		<span class='sf-icon-License list_svg'></span>
	</li>
	<li id='menuaddressbook' class='sidebar-item' on-click='setScreen("addressbook","menuaddressbook");'>
		<span class='sf-icon-About list_svg'></span>
	</li>
	<li id='menublockexplorer' class='sidebar-item' on-click='setScreen("blockexplorer","menublockexplorer");'>
		<span class='sf-icon-Notification list_svg'></span>
	</li>
	<li id='menusettings' class='sidebar-item' on-click='setScreen("settings","menusettings");'>
		<span class='sf-icon-Hardware list_svg'></span>
	</li>
	<li id='menucharts' class='sidebar-item' on-click='setScreen("charts","menucharts");'>
		<span class='sf-icon-Request list_svg'></span>
	</li>
</ul>`
}
