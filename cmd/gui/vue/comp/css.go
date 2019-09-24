package comp

var GetCoreCss = `

:root {
--blue: #3030cf;
--light-blue: #30cfcf;

--green: #30cf30;
--orange: #ff7500;
--yellow: #ffd500;
--red: #cf3030;
--purple: #cf30cf;

--dark: #303030;
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


.lineBottom{
	box-shadow: inset 0 -1px 0 rgba(0, 0, 0, .13);
}

.lineRight{
	box-shadow: inset -1px 0 0 rgba(0, 0, 0, .12);
}

.bgDark{
	background-color:var(--dark);
}
.bgLight{
	background-color:var(--light);
}
.bgMoreLight{
	background-color:var(--light-grayiii);
}


html, body{
	width:100%;
	height:100vh;
	margin:0;
	padding:0;
}


x.x{
	position:fixed;
	display: grid;
	grid-gap: 0;
	grid-template-columns: 60px 1fr;
	grid-template-rows: 60px 1fr;
	grid-template-areas:
	"Logo Header"
	"Sidebar Main"
	width:100%;
	height:100vh;
	margin:0;
	padding:0;
	top:0;
	left:0;
	right:0;
	bottom:0;
	overflow:hidden;
	background-color:var(--dark);
}


 .rwrap{
    position: relative;
    display: flex;
    width: 100%;
	height:100%;
	overflow:hidden;
	overflow-y:auto;
  }


.cwrap{
    position: relative;
    display: block;
    width: 100%;
	height:100%;
  }


.posRel{
    position: relative;
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

.grayGrad{
background:linear-gradient(#f4f4f4, #e4e4e4);
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

`
