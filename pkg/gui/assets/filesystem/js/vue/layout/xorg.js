var Xorg = new Vue({
	el: '#xorg', 
	name: 'Xorg',
	data () { return { 
	  rcvar }},
	template: `<div id="app" class="fullScreen bgDark flx lightTheme">
    <div id="display" class="fii">
      <div class="grid-container rwrap bgDark">
        <div class="flx fii Logo">
          <Logo />
        </div>
        <Header />
        <div class="Sidebar bgLight">
          <div class="Open"></div>
          <Nav />
          <div class="Side"></div>
        </div>
        <div id="main" class="grayGrad Main">
        <keep-alive><component v-bind:is="rcvar.isScreen"></component></keep-alive>
        </div>
      </div>
    </div>
  </div>`,
});