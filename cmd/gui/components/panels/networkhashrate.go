package panels

import "github.com/p9c/pod/cmd/gui"


func PanelNetworkHashRate()gui.DuOScomP{
    return gui.DuOScomP{
        Name:        "",
        ID:          "",
        Version:     "",
        Description: "",
        State:       "",
        Image:       "",
        URL:         "",
        CompType:    "",
        SubType:     "",
        Js:          "",
        Html:        "",
        Css:         "",
    }
}


Vue.component('PanelNetworkHashrate', {
    name: 'PanelNetworkHashrate',
    data () { return { 
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
                let value = rcvar.osNetHash;
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
},
		template: `<div class='posAbs rwrap'>
		<ejs-sparkline ref="networkHashrate" class="spark" id='networkHashrate-container' :height='this.height' :padding='padding' :width='this.width' :lineWidth='this.lineWidth' :type='this.type' :valueType='this.valueType' :fill='this.fill' :dataSource='this.dataSource' :axisSettings='this.axisSettings' :containerArea='this.containerArea' :border='this.border' xName='x' yName='yval'></ejs-sparkline>                        
	  <div style="color: #303030; font-size: 12px; position: absolute; top:12px; left: 15px;">Network hashrate</div>
	<div id="nhr" style="color: #d1a990;position: absolute; top: 25px; left: 15px;">R: 0H/s</div>
</div>`,
});