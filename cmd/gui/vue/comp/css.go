package comp

var GetCoreCss string = `

:root {
--blue: #3030cf;
--light-blue: #30cfcf;

--green: #30cf30;
--orange: #ff7500;
--yellow: #ffd500;
--red: #cf3030;
--purple: #cf30cf;

--black: #303030;
--white: #fcfcfc;

--dark-gray:#4d4d4d;
--gray:#808080;
--light-gray:#c0c0c0;

--light-green:#80cf80;
--light-orange:#cf9480;
--light-yellow:#cfcf80;
--light-red:#cf8080;
--light-purple:#cf80cf;

--dark-blue:#303080;
--dark-green:#308030;
--dark-orange:#804430;
--dark-yellow:#808030;
--dark-red:#803030;
--dark-purple:#803080;
  
--green-blue:#308080;
--light-green-blue:#80a8a8;
--dark-green-blue:#305858;
--green-orange:#80a830;
--light-green-orange:#a8bc80;
--dark-green-orange:#586c30;
--green-yellow:#80cf30;
--light-green-yellow:#a8cf80;
--dark-green-yellow:#588030;
--green-red:#808030;
--light-green-red:#a8a880;
--dark-green-red:#585830;
--blue-orange:#80a830;
--light-blue-orange:#a8bc80;
--dark-blue-orange:#583a58;
--bluered:#803080;
--light-blue-red:#a880a8;
--dark-blue-red:#583058;
--dark:#303030;
--dark-grayii:#424242;
--dark-grayi:#535353;
--dark-gray:#656565;
--gray:#808080;
--light-gray: #888888;
--light-grayi: #9a9a9a;
--light-grayii: #acacac;
--light-grayiii: #bdbdbd;
--light:#cfcfcf;











  --border-light: rgba(255, 255, 255, .62);
  --border-dark: rgba(0,0,0, .38);

  --trans-light: rgba(255, 255, 255, .24);
  --trans-dark: rgba(0,0,0, .24);
  --trans-gray: rgba(48,48,48, .38);

  --fonta:'Roboto';
  --fontb:'Abril Fatface';
  --fontc:'Oswald';
  --big-title:'Vollkorn SC';
  
  --base:var(--white);
  --pri: var(--blue);
  --sec: var(--light-blue);
  --btn-tx: var(--base);
  --btn-h-tx: #fff;
  --btn-bg: var(--dark-green-blue);
  --btn-h-bg: var(--green-blue);

  --space-02: .12rem;  /* 2px */
  --space-05: .25rem;  /* 4px */
  --space-1: .5rem;  /* 8px */
  --space-2: 1rem;   /* 16px */
  --space-3: 1.5rem; /* 24px */
  --space-4: 2rem;   /* 32px */
  --space-5: 2.5rem;   /* 40px */
  --space-6: 3rem;   /* 48px */
  --space-7: 3.5rem;   /* 56px */
  --space-8: 4rem;   /* 64px */

  

  --box-shadow-b: 0 1px 0 0 var(--black);
  --box-shadow-l: 0 1px 0 0 var(--white);
  --box-shadow-inset :inset 0 0 0 1px var(--sec);
}




html, body{
	width:100%;
	height:100vh;
	margin:0;
	padding:0;
}
#dev, #display, #container{
	display:flex;
	width:100%;
	height:100vh;
	margin:0;
	padding:0;
	overflow:hidden;
}

  .rwrap{
    position: relative;
    display: flex;
    width: 100%;
	height:100%;
	overflow:hidden;
	overflow-y:auto;
  }





.sidebar-content {
    padding: 14px;
    width: calc(98% - 60px);
}


  .dashboard-header {
    border-bottom: 1px solid rgba(0, 0, 0, 0.12);
    height: 59px;
    position: relative;
  }

  #analysisLayout.e-dashboardlayout.e-control .e-panel .e-panel-container .e-panel-header {
    border-bottom: 2px solid #e6e9ed !important;
    padding: 10px;
    height: 35px;
    margin: 0 15px 0 15px;
  }

  #analysisLayout.e-dashboardlayout  .e-panel-content {
    height: calc(100% - 35px) !important;
    overflow: hidden;
    width: 100%;
  }

  #sidebar-section {
    padding: 0px !important;
  }

  .e-bigger #search {
    display: none;
  }

  @media (max-width: 650px) {
    .e-bigger .searchContent {
      display: none;
    }
    .information{
        right:17% !important;
    }
  }

  @font-face {
    font-family: "e-sb-icons";
    src: url(data:application/x-font-ttf;charset=utf-8;base64,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cwBSAGUAZwB1AGwAYQByAGUALQBzAGIALQBpAGMAbwBuAHMAZQAtAHMAYgAtAGkAYwBvAG4AcwBWAGUAcgBzAGkAbwBuACAAMQAuADAAZQAtAHMAYgAtAGkAYwBvAG4AcwBGAG8AbgB0ACAAZwBlAG4AZQByAGEAdABlAGQAIAB1AHMAaQBuAGcAIABTAHkAbgBjAGYAdQBzAGkAbwBuACAATQBlAHQAcgBvACAAUwB0AHUAZABpAG8AdwB3AHcALgBzAHkAbgBjAGYAdQBzAGkAbwBuAC4AYwBvAG0AAAAAAgAAAAAAAAAKAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA2AQIBAwEEAQUBBgEHAQgBCQEKAQsBDAENAQ4BDwEQAREBEgETARQBFQEWARcBGAEZARoBGwEcAR0BHgEfASABIQEiASMBJAElASYBJwEoASkBKgErASwBLQEuAS8BMAExATIBMwE0ATUBNgE3AAZtaWMtMDMMYnVsbGV0cy0tLTA1CmJ1bGxldHMtd2YJd2FsbGV0LXdmCWhlbHAtLS0wMgloZWxwLS0tMDEQcGFpZC1pbnZvaWNlMy13ZhBwYWlkLWludm9pY2UyLXdmEHBhaWQtaW52b2ljZTEtd2YOZm9sZGVyLW9wZW4tMDEFdGFza3MKdGFzay0wMi13Zgd0YXNrLTAyB3Rhc2stMDEKYXJyb3d1cC13Zg1hcnJvdy1kb3duLXdmDWNhbGVuZGFyLS0wMS0OY2FsZW5kYXItMDItd2YOY2FsZW5kYXItMDMtd2YMYnVsbGV0cy0tLTAyDGJ1bGxldHMtLS0wMQpidWxsZXRzXzAxB2hvbWVfMDEIaG91c2UtMDcIaG91c2UtMDkIc2V0dGluZ3MNc2V0dGluZ3MtLS0xMQ1zZXR0aW5ncy0tLTEwC3NldHRpbmdzLTAyBXRpbWVyFXNob3J0LWFzY2VuZGluZy0wMS13Zgxwcm9maXQtMDEtd2YOc3RvY2staW5kZXgtdXAIY2hhcnQtd2YFZ3JhcGgIZ3JhcGgtMDEJZGF0YS1lZGl0EXRleHQtaGlnaGxpZ2h0LXdmDGRhdGEtZWRpdC13ZglzZWFyY2gtd2YLc2VhcmNoLWZpbmQLZmlsdGVyLWVkaXQOZmlsdGVyLS0tMDItd2YMZmlsdGVyLTEwLXdmB21lbnVfMDEKbWVudS0wMS13ZhFtZW51LWludGVyZmFjZS13Zgx1c2VyLXByb2ZpbGURdXNlci1wcm9maWxlLTEtd2YLdXNlci1hZGQtd2YEc2hvdwhjbG9jay13Zgh0aW1lci13ZgAA)
      format("truetype");
    font-weight: normal;
    font-style: normal;
  }

  [class^="sf-icon-"],
  [class*=" sf-icon-"] {
    font-family: "e-sb-icons" !important;
    speak: none;
    font-size: 55px;
    font-style: normal;
    font-weight: normal;
    font-variant: normal;
    text-transform: none;
    line-height: 1;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }

  #dashboardSidebar {
    text-align: center;
  }

  #dashboardSidebar .e-icons::before, #sidebar-section .e-icons::before {
    font-size: 20px;
  }

  .sidebar-item.current, .sidebar-item:hover {
    background: #efefef;
}

	.sidebar-item:hover{
	color:#fff;
	transition: all 900ms 0 ease;
	}
	.sidebar-item:current{
	transition: all 300ms 0 ease;
	}

  #menuoverview:hover {
    box-shadow: inset 60px 0 0 #30cf30;
  }
  #menuoverview.current {
    box-shadow: inset 4px 0 0 #30cf30;
	color:#30cf30;
  }


  #menufilter:hover {
    box-shadow: inset 60px 0 0 #30cfcf;
  }
  #menufilter.current {
    box-shadow: inset 4px 0 0 #30cfcf;
	color:#30cfcf;
  }


  #menusettings:hover {
    box-shadow: inset 60px 0 0 #cfcf30;
  }
  #menusettings.current {
    box-shadow: inset 4px 0 0 #cfcf30;
	color:#cfcf30;
  }

  #menucharts:hover {
    box-shadow: inset 60px 0 0 #3030cf;
  }
  #menucharts.current {
    box-shadow: inset 4px 0 0 #3030cf;
	color:#3030cf;
  }




  /* dockbar icon Style */

  .home::before {
    content: "\e718";
    font-family: "e-sb-icons";
  }

  .filter::before {
    content: "\e72a";
    font-family: "e-sb-icons";
  }

  .analyticsChart::before {
    content: "\e722";
    font-family: "e-sb-icons";
  }

  .analytics::before {
    content: "\e720";
    font-family: "e-sb-icons";
  }

  .session::before {
    content: "\e735";
    font-family: "e-sb-icons";
  }

  .profile::before {
    content: "\e730";
    font-family: "e-sb-icons";
  }

  .views::before {
    content: "\e733";
    font-family: "e-sb-icons";
  }

  .search::before {
    content: "\e728";
    font-family: "e-sb-icons";
    font-size: 14px;
    position: absolute;
    top: 5px;
    left: 12%;
    font-weight: 800;
  }

  #search:hover {
    background-color: transparent !important;
  }

  .settings::before {
    content: "\e71d";
    font-family: "e-sb-icons";
  }

  span.e-input-group.e-control-wrapper.e-ddl {
    left: 12%;
  }

  .expand::before,
  .expand::before {
    content: "\e72d";
    margin-left: 18px;
    font-family: "e-sb-icons";
    position: absolute;
    top: 12%;
  }

  .right-content {
    float: right;
    
    
  }

  #right-sidebar {
    display: none;
  }

  .e-dock.e-close span.e-text {
    display: none;
  }

  .e-dock.e-open span.e-text {
    display: inline-block;
  }

  #dashboardSidebar li {
    list-style-type: none;
    cursor: pointer;
    padding: 5px;
  }

  #dashboardSidebar ul {
    padding: 0px;
  }

  span.e-icons {
    line-height: 2;
  }

  .e-open .e-icons {
    margin-right: 16px;
  }

  .analysis {
    font-size: 18px;
    padding: 12px;
    text-align: left;
    vertical-align: middle;
  }

  #search {
    margin-left: 10px;
    text-indent: 8px;
  }

  .searchContent .e-input-group.e-control-wrapper.e-ddl {
    height: 28px !important;
  }


  .card-content.text {
    font-size: 15px;
    text-align: right;
    color: #66696b;
  }

  .card {
    margin-right: 5%;
    margin-top: 10%;
  }

  .card-content.number {
    font-size: 16px;
    text-align: right;
    padding-top: 10px;
  }

  #header-avatar.e-avatar.image {
    background-image: url('./images/pic01.png');
    background-repeat: no-repeat;
    background-size: cover;
    background-position: center;
  }

  #sidebarTarget {
    background: linear-gradient(-141deg, #eaeaea 14%, #dadada 100%);
  }

  .markerTemplate {
    font-size: 12px;
    color: white;
    text-shadow: 0px 1px 1px black;
    font-weight: 500;
  }

  .logoContainer {
    display: inline-block;
    width: 59px;
    height: 60px;
    border-right: 1px solid rgba(0, 0, 0, 0.12);
	background-color:#303030;
  }

  .searchContent {
    display: inline-block;
    position: absolute;
    left: 60px;
  }


  .e-dock {
    padding-top: 8px;
  }

  .right-content div {
    display: inline-block;
  }

  .information {
    font-size: 12px;
  }

  .information span {
    float: left;
  }

  .text-content {   
    font-size: 17px;    
    margin-left: 10px;
    margin-top: 10px;
  }

  .card .e-icons {
    position: absolute;
    top: 20%;
    left: 12%;
    width: 60px;
    height: 60px;
    text-align: center;
    border: 1px solid;
    line-height: 60px;
    color: #cfcfcf;
    background: #303030;
    border-radius: 60px;
  }

  .card .home::before {
    font-size: 25px;
  }

  .dashboardParent , .chart-content, .map-content, .maps-content{
    height: 100%;
    width:100%;
  }

   /* styles for highcontrast theme */



  body.highcontrast #dashboardSidebar .e-icons::before {
    color: #cfcfcf;
  }

  body.highcontrast #dashboardSidebar li.sidebar-item.filterHover,
  body.highcontrast #dashboardSidebar {
    background: #303030;
  }

body.highcontrast #analysisLayout.e-dashboardlayout.e-control .e-panel {
    background: #303030;
}

body.highcontrast #analysisLayout.e-dashboardlayout.e-control .e-panel .e-panel-container .e-panel-header {
    color: rgba(255, 255, 255, 0.54);
}






.e-panel-header{
height:48px;
}
.e-panel-header h3{
margin:0;
}
.card-value {
    font-family: 'Roboto';
    font-weight: 400;
    font-size: 32px;
    color: #727272;
    margin-bottom: 0px;
    margin-top: 0px;
    letter-spacing: 0;
}
.card-text {
    font-family: 'Roboto';
    font-weight: 400;
    font-size: 18px;
    letter-spacing: 0;
    color: #303030;
}



#layoutconfig-panel{
	text-align:right;
}
.lytconf{
display:flex;
flex-direction:column;
justify-items:right;
}
.lytconf button{
	margin:0;
}
input{
flex:1;
background:#fff!important;
}
.noPanelFrame{
	box-shadow:0 0 0 #fff!important;
	border: none!important;
}
.e-panel.noPanelFrame{
background:transparent!important;
}
.noPanelFrame .e-panel-header{
	display:none;
}

.noPanelFrame .e-panel-header{
	display:none;
}
.e-panel-content.noPanelFrame{
height:100%!important;
display:flex;
}
.noPanelFrame .e-float-input, .noPanelFrame .e-input-group{
margin:0;
}
.noPanelFrame button{
margin:0;
}
.flx{
	display:flex;
}
.flc{
	flex-direction: column;
}
.fii{
	flex:1;
}

.bgfff{
background-color:#fff;
}


.marginTop{
margin-top:30px;
}

.justifyCenter{
justify-content: center;
}
.justifyBetween{
justify-content: space-between;
}

.justifyEvenly{
justify-content: space-evenly;
}

.itemsCenter{
	align-items: center;
}

.noMargin{
	margin:0!important;
}

.noPadding{
	padding:0!important;
}


#main.lightTheme{
background:#cfcfcf;
color:#303030;
}
#main{
background:#303030;
color:#cfcfcf;
}



.lightTheme .logofill{
background:#cfcfcf;
fill:#cfcfcf;
}
.logofill{
background:#303030;
fill:#303030;
}



.lightTheme .e-sidebar {
	background:#cfcfcf;
	color:#303030;
}
.e-sidebar {
	background:#303030;
	color:#cfcfcf;
}


.lightTheme .e-panel-header{
background:#cfcfcf;
}

.lightTheme .e-panel-header h3{
color:#303030;
}

.e-panel-header{
background:#303030;
}
.e-panel-header h3{
color:#cfcfcf;
}
.e-float-input{
background-color:#fff!important;
}

.e-float-input input{
flex:1;
}


#main.lightTheme .menu{
background:#303030;
}
.menu{
background:#cfcfcf;
}
#main.lightTheme .svgColor{
stroke:#cfcfcf;
}
.svgColor{
stroke:#303030;
}





.hide{
display:none;
}

.logo{
justify-content:center;
align-items:center;
}
.logo svg{
width:auto;
height:48px;
}



.form-group{
	display:flex;
	margin-bottom:8px;
	padding-bottom:8px;
	justify-content:space-between;
	border-bottom:1px solid #cfcfcf;
}
































@font-face {
  font-family: "button-icons";
  src: url(data:application/x-font-ttf;charset=utf-8;base64,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)
    format("truetype");
  font-weight: normal;
  font-style: normal;
}

#button-control .e-btn-sb-icons {
  font-family: "button-icons";
  line-height: 1;
  font-style: normal;
  font-weight: normal;
  font-variant: normal;
  text-transform: none;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#button-control .e-play-icon::before {
  content: "\e701";
}

#button-control .e-pause-icon::before {
  content: "\e705";
}


`
