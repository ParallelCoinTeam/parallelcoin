package pnl

import "github.com/p9c/pod/cmd/gui/mod"

func WalletStatus() mod.DuOScomp {
	return mod.DuOScomp{
		IsApp:    true,
		Name:     "Wallet status",
		ID:       "panelwalletstatus",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "status",
		Js: `
			data () { return { 
			duoSystem }}, 
		`,
		Template: `<div class="rwrap"><div tabindex="0" class="e-card flx flc justifyBetween duoCard"><span v-html="this.duoSystem.status"></span></div></div>`,
		Css: `

.e-card .e-card-header .e-card-header-image.balance {
	position:absolute;
	bottom:0;
	left:0;
    background-image: url('./football.png');
}

.e-card.duoCard {
width: 100%;
height:100%;
background: linear-gradient(-135deg, #f4f4f4 0%, #d4d4d4) 100%;
}
.e-card.duoCard div.e-card-sub-title{
text-align:right;

}

.e-card.duoCard  .e-card-sub-title span{
	font-size:36px;
	font-weight:100;
	letter-spacing:-1px;
}
.e-card.duoCard .e-card-content{
margin-left:auto;
}
		`,
	}
}

func Status() mod.DuOScomp {
	return mod.DuOScomp{
		IsApp:    true,
		Name:     "Status",
		ID:       "panelstatus",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "status",
		Js: `
			data () { return { 
			duoSystem }}, 
		`,
		Template: `<div class="rwrap">


    <ul class="rf flx flc noPadding justifyEvenly">
        <li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Version: </span><strong class="rcx6"><span v-html="this.duoSystem.status.ver"></span></strong></li>
		<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Wallet version: </span><strong class="rcx6"><span v-html="this.duoSystem.status.walletver.podjsonrpcapi.versionstring"></span></strong></li>

        <li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Uptime: </span><strong class="rcx6"><span v-html="this.duoSystem.status.uptime"></span></strong></li>


		<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Memory: </span><strong class="rcx6"><span v-html="this.duoSystem.status.mem.total"></span></strong></li>
		<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Disk: </span><strong class="rcx6"><span v-html="this.duoSystem.status.disk.total"></span></strong></li>

		<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Chain: </span><strong class="rcx6"><span v-html="this.duoSystem.status.net"></span></strong></li>
		<li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Blocks: </span><strong class="rcx6"><span v-html="this.duoSystem.status.blockcount"></span></strong></li>
        <li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Connections: </span><strong class="rcx6"><span v-html="this.duoSystem.status.connectioncount"></span></strong></li>
    </ul>

</div>`,
		Css: `



		`,
	}
}

func LocalHashRate() mod.DuOScomp {
	return mod.DuOScomp{
		IsApp:    true,
		Name:     "Local Hashrate",
		ID:       "panellocalhashrate",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "status",
		Js: `
   data:function(){
			return { 
			duoSystem,
        height: '100%',
        width: '100%',
    	padding: { left: 0, right: 0, bottom: 0, top: 0},
        axisSettings: {
            minY: 0, maxY: 99999
        },
        containerArea: {
            background: 'white',
            border: {
                color: '#dcdfe0',
                width: 0
            }
        },
        border: {
            color: '#0358a0',
            width: 1
        },
        fill: '#e8f2fc',
        type: 'Area',
        valueType: 'Numeric',
        dataSource:[
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 }
],
		lineWidth: 1
    }
},
mounted(){
    this.update();
},
methods:{
    update: function() {
        let spark = document.getElementById('localHashrate-container');
        let gauge = this.$refs.localHashrate.ej2Instances;
        let temp = gauge.dataSource.length - 1;
        this.update = setInterval(function() {
            if (gauge.element.className.indexOf('e-sparkline') > -1) {
                let value = duoSystem.status.hashrate;
                gauge.dataSource.push({ x: ++temp, yval: value });
                gauge.dataSource.shift();
                gauge.refresh();
                let net = document.getElementById('lhr');
                if (net) {
                net.innerHTML = 'R: ' + value.toFixed(0) + 'H/s';
                }
            }
        }, 500);
    }
}
		`,
		Template: `<div class='rwrap'>
                        <ejs-sparkline ref="localHashrate" class="spark" id='localHashrate-container' :height='this.height' :width='this.width' :padding='padding' :lineWidth='this.lineWidth' :type='this.type' :valueType='this.valueType' :fill='this.fill' :dataSource='this.dataSource' :axisSettings='this.axisSettings' :containerArea='this.containerArea' :border='this.border' xName='x' yName='yval'></ejs-sparkline>                        
                      <div style="color: #000000; font-size: 12px; position: absolute; top:12px; left: 15px;">
                        <b>Local hashrate</b>
                    </div>
                    <div id="lhr" style="color: #d1a990;position: absolute; top:25px; left: 15px;">R: 0H/s</div>
</div>`,
		Css: `

  .spark{
    height: 100%;
    width:100%;
  }
		`,
	}
}
func NetworkHashRate() mod.DuOScomp {
	return mod.DuOScomp{
		IsApp:    true,
		Name:     "Network Hashrate",
		ID:       "panelnetworkhashrate",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "status",
		Js: `
   data:function(){
			return { 
			duoSystem,
        height: '100%',
        width: '100%',
    	padding: { left: 0, right: 0, bottom: 0, top: 0},
        axisSettings: {
            minY: 0, maxY: 9999999
        },
        containerArea: {
            background: 'white',
            border: {
                color: '#dcdfe0',
                width: 0
            }
        },
        border: {
            color: '#cf8030',
            width: 1
        },
        fill: '#cfa880',
        type: 'Area',
        valueType: 'Numeric',
        dataSource:[
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 },
    { x: 0, yval: 0 }
],
		lineWidth: 1
    }
},
mounted(){
    this.update();
},
methods:{
    update: function() {
        let spark = document.getElementById('networkHashrate-container');
        let gauge = this.$refs.networkHashrate.ej2Instances;
        let temp = gauge.dataSource.length - 1;
        this.update = setInterval(function() {
            if (gauge.element.className.indexOf('e-sparkline') > -1) {
                let value = duoSystem.status.networkhashrate;
                gauge.dataSource.push({ x: ++temp, yval: value });
                gauge.dataSource.shift();
                gauge.refresh();
                let net = document.getElementById('nhr');
                if (net) {
                net.innerHTML = 'R: ' + value.toFixed(0) + 'H/s';
                }
            }
        }, 500);
    }
}
		`,
		Template: `<div class='rwrap'>
                        <ejs-sparkline ref="networkHashrate" class="spark" id='networkHashrate-container' :height='this.height' :padding='padding' :width='this.width' :lineWidth='this.lineWidth' :type='this.type' :valueType='this.valueType' :fill='this.fill' :dataSource='this.dataSource' :axisSettings='this.axisSettings' :containerArea='this.containerArea' :border='this.border' xName='x' yName='yval'></ejs-sparkline>                        
                      <div style="color: #303030; font-size: 12px; position: absolute; top:12px; left: 15px;">Network hashrate</div>
                    <div id="nhr" style="color: #d1a990;position: absolute; top: 25px; left: 15px;">R: 0H/s</div>
</div>`,
		Css: `

  .spark{
    height: 100%;
    width:100%;
  }
		`,
	}
}
