package sys

import "github.com/parallelcointeam/parallelcoin/cmd/gui/vue/mod"

func Screen() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Screen",
		ID:       "screen",
		Version:  "0.0.0",
		CompType: "core",
		SubType:  "screen",
		Js: `
	data () {
		return { 
			duoSystem,
       		spacing: [10,10],
          	header:'Add a Content',
          	target:'.control-section',
          	width:'43%',
          	showCloseIcon: true,
          	contenttemplateVue:'<div id="dialogcontent"><div><div id="linetemplate"><p class="dialog-text">Linechart (1x1) </p></div><div id="pietemplate"><p class="dialog-text">Piechart (1x1) </p></div><div id="splinetemplate"><p class="dialog-text">Splinechart (2x1) </p></div></div></div><div id="headerTemplate"><span id="close" class="e-template-icon e-clear-icon"></span></div></div>',
		}},
		methods:{
            onPanelResize: function(args) {
            	var dashboardObject = this.$refs.DuoDashboard;
        		if (dashboardObject && args.element && args.element.querySelector('.e-panel-container .e-panel-content div') && dashboardObject.$el.querySelector('.e-holder')) {
            		var chartObj = (args.element.querySelector('.e-panel-container .e-panel-content div.e-control')).ej2_instances[0];
            		var holderElementHeight = parseInt((dashboardObject.$el.querySelector('.e-holder')).style.height, 10);
            		var panelContanierElement = args.element.querySelector('.e-panel-content');
            		panelContanierElement.style.height = holderElementHeight - 35 + 'px';
            		chartObj.height = '95%';
            		chartObj.width = '100%';
            		chartObj.refresh();
				}
			},

}
			`,
		Template: `
<div>

    <ejs-dashboardlayout ref="DuoDashboard"
:columns="9"
:allowResizing="false"
:allowDragging="false"
:cellSpacing="false"
:resizeStop="onPanelResize">
		<e-panels>
			<e-panel v-for="panel in this.duoSystem.screenPanels" :id="panel.id" :sizeX="panel.sizeX" :sizeY="panel.sizeY" :row="panel.row" :col="panel.col" :header="panel.header" :content="panel.content" class="rwrap""></e-panel>
		</e-panels>
	</ejs-dashboardlayout>   
</ul>

	<ejs-dialog :header='header' ref="dialogObj" :content='contenttemplateVue' :showCloseIcon='showCloseIcon' :target='target' :width='width' :visible='false' :isModal='true'></ejs-dialog>
</div>
		`,
		Css: `

		`,
	}
}
