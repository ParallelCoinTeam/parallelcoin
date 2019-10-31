<script>
export default {
	name: 'PanelTimeBalance',
  data: function() {
    return {
		duoSystem,
      	theme: "Material",
      primaryXAxis: {
		valueType: "DateTime",
        intervalType: "Auto",
        edgeLabelPlacement: "Shift",
        majorGridLines: { width: 0 }
      },
      primaryYAxis: {
        labelFormat: "{value} DUO",
        rangePadding: "None",
        minimum: 0,
        maximum: system.data.d.txsex.balanceheight,
        interval: 20,
        lineStyle: { width: 0 },
        majorTickLines: { width: 0 },
        minorTickLines: { width: 0 }
      },
      chartArea: {
        border: {
          width: 0
        }
      },
      marker: {
        visible: true,
        height: 10,
        width: 10
      },
      tooltip: {
        enable: true
      },
      title: "Inflation - Consumer Price"
    	}
  }
}
</script>
		
<template>
<div class="rwrap">
  <div class="control-section">
    <div align='center'>
        <ejs-chart style='display:block' :theme='this.theme' align='center' id='chartcontainer' :title='this.title' :primaryXAxis='this.primaryXAxis' :primaryYAxis='this.primaryYAxis'
            :tooltip='this.tooltip' :chartArea='this.chartArea' width='100%' height='100%'>
            <e-series-collection>
                <e-series :dataSource='system.data.d.txsex.txs' type='Line' xName='time' yName='balance' name='Germany' width=2 :marker='this.marker'> </e-series>
            </e-series-collection>
        </ejs-chart>
    </div>
</div>
</div>
		</template>