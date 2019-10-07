package css

func CSS(root, grid, nav string) string {
	return `<style amp-custom>
` + root +
		grid + nav + `
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
.bgfff{
background-color:#fff!important;
}

.grayGrad{
background:linear-gradient(#f4f4f4, #e4e4e4);
}
html, body{
width:100%;
height:100vh;
margin:0;
padding:0;
}
#x{
position:fixed;
margin:0;
padding:0;
top:0;
left:0;
right:0;
bottom:0;
overflow:hidden;
background-color:var(--dark);
}
display#display{
position:relative;
display:block;
width:100%;
height:100%;
}

.rwrap{
position: relative;
display: block;
width: 100%;
height:100%;
overflow:hidden;
}
.cwrap{
position: relative;
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


.marZ{
margin:0!important;
}
.padZ{
padding:0!important;
}


.lsn{
list-style:none!important;
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


.baseMargin{
margin:15!important;
}






.hide{
display:none;
}

.Logo{
justify-content:center;
align-items:center;
background:var(--dark);
}
.Logo svg{
width:auto;
height:48px;
}



.logofill{
fill:var(--light);
background:var(--light);
}





</style>`
}
