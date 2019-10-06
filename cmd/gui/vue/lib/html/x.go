package html

func VUEx(logo, header, nav, overview string) string {
	return `<x id="x" v-show="!this.duOSys.bios.isBoot" class="bgDark lightTheme"><display id="display">
<div class="grid-container grayGrad">
	<div class="flx fii Logo">` + logo + `</div>
	<div class="Header bgLight">` + header + `</div>
	<div class="Sidebar bgLight">
		<div class="Open"></div>
		<div class="Nav">` + nav + `</div>
		<div class="Side"></div>
	</div>
	<main id="main" class="Main">` + overview + `</main>
</div>
</display></x>`
}
