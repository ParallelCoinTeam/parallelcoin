package html

func VUEx(logo, header, nav, overview string) string {
	return `<div id="x" v-show="!this.duOSys.bios.isBoot" class="bgDark lightTheme"><div id="display">
<div class="grid-container bgDark">
	<div class="flx fii Logo">` + logo + `</div>
	<div class="Header bgLight">` + header + `</div>
	<div class="Sidebar bgLight">
		<div class="Open"></div>
		<div class="Nav">` + nav + `</div>
		<div class="Side"></div>
	</div>
	<div id="main" class="grayGrad Main">` + overview + `</div>
</div>
</div></div>`
}
