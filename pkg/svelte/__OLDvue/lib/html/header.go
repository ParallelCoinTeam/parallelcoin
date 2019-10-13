package html

func VUEheader() string {
	return `
<div class="h1"></div>
<div class="h2"></div>
<div class="h3">
	<div class="searchContent">
		<div class="analysis">ParallelCoin</div>
	</div>
</div>
<div class="h4">
	<button ref="toggleFullScreenBtn">FullScreen</button>
</div>
<div class="h5"></div>
<div class="h6"></div>
<div class="h7"></div>
<div class="h8">
	<button id="toggle" ref="toggleBoardbtn" class="e-btn e-info"  cssClass="e-flat" iconCss="e-icons burg-icon" isToggle="true" v-on:click.native="btnClick">Open</button>
</div>`
}
