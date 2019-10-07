package sys

import "github.com/p9c/pod/cmd/gui/vue/mod"

func Nav() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Nav",
		ID:       "nav",
		Version:  "0.0.1",
		CompType: "core",
		SubType:  "nav",
		Js: `
		methods: { toggle: function(){ this.navOpen = !this.navOpen }}, 
		`,
		Template: `
		<div class="nav">
			<div class="pie pie1" @click="document.body.classList.remove('active')">
				<div class="pie-color pie-color1">1</div>
			</div>
			<div class="pie pie2" @click="document.body.classList.remove('active')">
				<div class="pie-color pie-color2">2</div>
			</div>
			<div class="pie pie3" @click="this.duoSystem.isDev = !duoSystem.isDev; document.body.classList.remove('active');">
				<div class="pie-color pie-color3">3</div>
			</div>
			<div class="menu" @click="document.body.classList.toggle('active');">
				<svg class="hamburger" xmlns="http://www.w3.org/2000/svg" width="100" height="100" viewBox="0 0 100 100"><g fill="none" class="svgColor" stroke-width="7.999" stroke-linecap="round"><path d="M 55,26.000284 L 24.056276,25.999716" /><path d="M 24.056276,49.999716 L 75.943724,50.000284" /><path d="M 45,73.999716 L 75.943724,74.000284" /><path d="M 75.943724,26.000284 L 45,25.999716" /><path d="M 24.056276,73.999716 L 55,74.000284" /></g></svg>
			</div>
		</div>
		`,
		Css: `
		.nav{
			position: fixed;
			width: 100px; 
			height: 100px;
			left: 0;
			top: 0;
			z-index: 11111;
			cursor:pointer;
		  }
		.icon{
			width: 32px; 
			height:32px;
		 }
		  .tjar{
			 background: #30cfccf;
			 fill:#30cfccf;
		   }
		 .pie {
		   background: #000;
		   border-radius: 50%;
		   box-shadow: 0 0 4px 5px rgba(0, 0, 0, 0.2);
		   cursor: pointer;
		   height: 400px;
		   left: -200px;
		   position: absolute;
		   top: -200px;
		   width: 400px;
		   transform: translateX(-200px) translateY(-200px);
		   transition: transform 300ms;
		 }
		 .pie-color:hover {
		   opacity: 0.85;
		 }
		 .pie-color:active {
		   opacity: 0.7;
		 }
		 .pie1 {
		   clip-path: polygon(200px 200px, 344px 450px, 0 450px);
		   transition-delay: 30ms;
		 }
		 .pie2 {
		   clip-path: polygon(200px 200px, 344px 450px, 450px 344px);
		   transition-delay: 60ms;
		 }
		 .pie3 {
		   clip-path: polygon(200px 200px, 450px 0, 450px 344px);
		   transition-delay: 90ms;
		 }
		 .pie-color {
		   width: 100%;
		   height: 100%;
		   border-radius: 50%;
		 }
		 .pie-color1 {
		   background: linear-gradient(135deg, #f12711, #f5af19);
		   clip-path: polygon(200px 200px, 344px 450px, 0 450px);
		 }
		 .pie-color2 {
		   background: linear-gradient(135deg, #444, #7e84f9);
		   clip-path: polygon(200px 200px, 344px 450px, 450px 344px);
		 }
		 .pie-color3 {
		   background: linear-gradient(135deg, #444, #b7e13f);
		   clip-path: polygon(200px 200px, 450px 0, 450px 344px);
		 }
		 .card {
		   left: 216px;
		   position: absolute;
		   top: 300px;
		   width: 46px;
		 }
		 .discount {
		   left: 288px;
		   position: absolute;
		   top: 258px;
		   width: 46px;
		 }
		 .cart {
		   left: 324px;
		   position: absolute;
		   top: 188px;
		   width: 46px;
		 }
		 .menu {
		   border-radius: 50%;
		   box-shadow: 0 0 4px 5px rgba(0, 0, 0, 0.2);
		   cursor: pointer;
		   height: 200px;
		   left: -100px;
		   position: absolute;
		   top: -100px;
		   width: 200px;
		 }
		 .hamburger {
		   cursor: pointer;
		   height: 46px;
		   left: 58%;
		   position: relative;
		   top: 58%;
		   width: 46px;
		 }
		 .title {
		   color: white;
		   font-family: "Crimson Text", serif;
		   font-size: 80px;
		   line-height: 84px;
		   margin-top: 60px;
		   text-align: center;
		 }
		 .body {
		   color: white;
		   font-family: "Work Sans", sans-serif;
		   font-size: 20px;
		   justify-content: center;
		   line-height: 28px;
		   margin: 30px auto;
		   max-width: 600px;
		   text-align: center;
		 }
		 .hamburger path {
		   transition: transform 300ms;
		 }
		 .hamburger path:nth-child(1) {
		   transform-origin: 25% 29%;
		 }
		 .hamburger path:nth-child(2) {
		   transform-origin: 50% 50%;
		 }
		 .hamburger path:nth-child(3) {
		   transform-origin: 75% 72%;
		 }
		 .hamburger path:nth-child(4) {
		   transform-origin: 75% 29%;
		 }
		 .hamburger path:nth-child(5) {
		   transform-origin: 25% 72%;
		 }
		 .active .pie {
		   transform: translateX(0) translateY(0);
		 }
		 .active .hamburger path:nth-child(1) {
		   transform: rotate(45deg);
		 }
		 .active .hamburger path:nth-child(2) {
		   transform: scaleX(0);
		 }
		 .active .hamburger path:nth-child(3) {
		   transform: rotate(45deg);
		 }
		 .active .hamburger path:nth-child(4) {
		   transform: rotate(-45deg);
		 }
		 .active .hamburger path:nth-child(5) {
		   transform: rotate(-45deg);
		 }	 
		`,
	}
}
