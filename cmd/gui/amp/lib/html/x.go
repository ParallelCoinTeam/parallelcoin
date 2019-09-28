package html

func AMPx(logo, header, nav string) string {
	return `<x id="x" v-show="!this.duoSystem.bios.isBoot" class="bgDark lightTheme"><display id="display">
<div class="grid-container grayGrad">
	<div class="flx fii Logo">` + logo + `</div>
	<div class="Header bgLight">` + header + `</div>
	<div class="Sidebar bgLight">
		<div class="Open"></div>
		<div class="Nav">` + nav + `</div>
		<div class="Side"></div>
	</div>
	<div class="Main">
		<div class="Balance"></div>
		<div class="Address"></div>
		<div class="Amount"></div>
		<div class="NetHR"></div>
		<div class="LocalHR"></div>
		<div class="BottomLeft">
			<div class="b1"></div>
			<div class="b2"></div>
			<div class="b3"></div>
			<div class="b4"></div>
			<div class="b5"></div>
			<div class="b6"></div>
		</div>
		<div class="Bottom"></div>
		<div class="Log"></div>
		<div class="Status"></div>
		<div class="Mid"></div>
		<div class="Txs"></div>
	</div>
</div>
</display></x>`
}
