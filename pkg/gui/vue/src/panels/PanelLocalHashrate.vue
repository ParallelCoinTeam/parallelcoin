<script>
import Vue from "vue";

import { SparklinePlugin } from '@syncfusion/ej2-vue-charts';
import { duoSystem } from '../DuOSyS.js'

Vue.use(SparklinePlugin);

export default {
   data:function(){
			return { 
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
        // let spark = document.getElementById('localHashrate-container');
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
}
		</script>

<template>
<div class='rwrap'>
                        <ejs-sparkline ref="localHashrate" class="spark" id='localHashrate-container' :height='this.height' :width='this.width' :padding='padding' :lineWidth='this.lineWidth' :type='this.type' :valueType='this.valueType' :fill='this.fill' :dataSource='this.dataSource' :axisSettings='this.axisSettings' :containerArea='this.containerArea' :border='this.border' xName='x' yName='yval'></ejs-sparkline>                        
                      <div style="color: #000000; font-size: 12px; position: absolute; top:12px; left: 15px;">
                        <b>Local hashrate</b>
                    </div>
                    <div id="lhr" style="color: #d1a990;position: absolute; top:25px; left: 15px;">R: 0H/s</div>
</div>
</template>
<style>
  .spark{
    height: 100%;
    width:100%;
  }
</style>