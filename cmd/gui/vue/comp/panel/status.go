package panel

import "github.com/parallelcointeam/parallelcoin/cmd/gui/vue/mod"

func WalletStatus() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
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
		Template: `<div class="rwrap">

                           
    <ul class="rf flx flc noPadding justifyEvenly">
        <li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Balance: </span><strong class="rcx6"><span v-html="this.duoSystem.status.balance.balance"></span></strong></li>
        <li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Unconfirmed: </span><strong class="rcx6"><span v-html="this.duoSystem.status.balance.unconfirmed"></span></strong></li>
        <li class="flx fwd spb htg rr"><span class="rcx2"></span><span class="rcx4">Transactions: </span><strong class="rcx6"><span v-html="this.duoSystem.status.txsnumber"></span></strong></li>
    </ul>
</div>`,
		Css: `



		`,
	}
}

func Status() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
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

func LocalHashRate() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
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
            minimum: 0,
            maximum: 120000,
            radius: '100%',
            pointerRadius: '100%',
            margin: { left: 0, right: 0, top: 0, bottom: 0 },
            lineStyle: { width: 0 },
            majorTicks: { width: 0, },
            minorTicks: { width: 0 },
            pointerWidth: 7,
            labelStyle: { useRangeColor: false, position: 'Outside', autoAngle: true,
            font: { size: '12px', fontFamily: 'Roboto' } },
            startAngle: 270, 
            endAngle: 90,
            color: '#757575',
            animation: { enable: true, duration: 900 },
            cap: {
                    radius: 8,
                    color: '#757575',
                    border: { width: 0 }
                },
            needleTail: {
                    color: '#757575',
                    length: '15%'
            },

            annotations: [
                {
                    content: '<div id="templateWrapLocal"><div class="des"><div id="pointerannotationLocal" style="width:90px;text-align:center;font-size:8px;font-family:Roboto">${pointers[0].value} Hash/second</div></div></div>',
                    angle: 0, zIndex: '1',
                    radius: '30%'
                }
            ],
            ranges: [
                {
                    start: 0,
                    end: 20000,
                    startWidth: 5, endWidth: 10,
                    radius: '100%',
                    color: '#cf3030',
                },
                {
                    start: 20000,
                    end: 40000,
                    startWidth: 10, endWidth: 15,
                    radius: '100%',
                    color: '#cf8030',
                }, {
                    start: 40000,
                    end: 60000,
                    startWidth: 15, endWidth: 20,
                    radius: '100%',
                    color: '#cfa880',
                },
                {
                    start: 60000,
                    end: 80000,
                    startWidth: 20, endWidth: 25,
                    radius: '100%',
                    color: '#80cf30',
                },
                {
                    start: 80000,
                    end: 100000,
                    startWidth: 25, endWidth: 30,
                    radius: '100%',
                    color: '#30CF30',
                },
                {
                    start: 100000,
                    end: 120000,
                    startWidth: 30, endWidth: 35,
                    radius: '100%',
                    color: '#588030',
                }
            ]
    }
},
		`,
		Template: `<div class='rwrap'>
<ejs-circulargauge ref="localHashRate" style='display:block;height:100%;width:100%;' align='center' width='100%' height='100%' id='localHashRate-container' :margin='margin' moveToCenter='true'>
<e-axes>
    <e-axis class="chart-content" :startAngle='startAngle' :endAngle='endAngle' :lineStyle='lineStyle' :labelStyle='labelStyle' :majorTicks='majorTicks' :minorTicks='minorTicks' 
    :radius='radius' :minimum='minimum' :maximum='maximum' :annotations='annotations' :ranges='ranges'>
      <e-pointers>
          <e-pointer :value='this.duoSystem.status.hashrate' :radius='pointerRadius' :color='color' :pointerWidth='pointerWidth' :animation='animation' :cap='cap' :needleTail="needleTail"></e-pointer>
      </e-pointers>  
    </e-axis>
</e-axes>
</ejs-circulargauge>
</div>`,
		Css: `

  #localHashRate{
    height: 100%;
    width:100%;
  }
		`,
	}
}

func NetworkHashRate() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
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
            minimum: 0,
            maximum: 12000000,
            radius: '100%',
            pointerRadius: '100%',
            margin: { left: 0, right: 0, top: 0, bottom: 0 },
            lineStyle: { width: 0 },
            majorTicks: { width: 0, },
            minorTicks: { width: 0 },
            pointerWidth: 7,
            labelStyle: { useRangeColor: false, position: 'Outside', autoAngle: true,
            font: { size: '8px', fontFamily: 'Roboto' } },
            startAngle: 270, 
            endAngle: 90,
            color: '#757575',
            animation: { enable: true, duration: 900 },
            cap: {
                    radius: 8,
                    color: '#757575',
                    border: { width: 0 }
                },
            needleTail: {
                    color: '#757575',
                    length: '15%'
            },

            annotations: [
                {
                    content: '<div id="templateWrapNetwork"><div class="des"><div id="pointerannotationNetwork" style="width:90px;text-align:center;font-size:12px;font-family:Roboto">${pointers[0].value} Hash/second</div></div></div>',
                    angle: 0, zIndex: '1',
                    radius: '30%'
                }
            ],
            ranges: [
                {
                    start: 0,
                    end: 2000000,
                    startWidth: 5, endWidth: 10,
                    radius: '102%',
                    color: '#cf3030',
                },
                {
                    start: 2000000,
                    end: 4000000,
                    startWidth: 10, endWidth: 15,
                    radius: '102%',
                    color: '#cf8030',
                }, {
                    start: 4000000,
                    end: 6000000,
                    startWidth: 15, endWidth: 20,
                    radius: '102%',
                    color: '#cfa880',
                },
                {
                    start: 6000000,
                    end: 8000000,
                    startWidth: 20, endWidth: 25,
                    radius: '102%',
                    color: '#80cf30',
                },
                {
                    start: 8000000,
                    end: 10000000,
                    startWidth: 25, endWidth: 30,
                    radius: '102%',
                    color: '#30CF30',
                },
                {
                    start: 10000000,
                    end: 12000000,
                    startWidth: 30, endWidth: 35,
                    radius: '102%',
                    color: '#588030',
                }
            ]
    }
},
		`,
		Template: `<div class='rwrap'>
<ejs-circulargauge ref="networkHashRate" style='display:block;height:100%;width:100%;' align='center' width='100%' height='100%' id='networkHashRate-container' :margin='margin' moveToCenter='true'>
<e-axes>
    <e-axis class="chart-content" :startAngle='startAngle' :endAngle='endAngle' :lineStyle='lineStyle' :labelStyle='labelStyle' :majorTicks='majorTicks' :minorTicks='minorTicks' 
    :radius='radius' :minimum='minimum' :maximum='maximum' :annotations='annotations' :ranges='ranges'>
      <e-pointers>
          <e-pointer :value='this.duoSystem.status.networkhashrate' :radius='pointerRadius' :color='color' :pointerWidth='pointerWidth' :animation='animation' :cap='cap' :needleTail="needleTail"></e-pointer>
      </e-pointers>  
    </e-axis>
</e-axes>
</ejs-circulargauge>
</div>`,
		Css: `

  #networkHashRate{
    height: 100%;
    width:100%;
  }
		`,
	}
}
