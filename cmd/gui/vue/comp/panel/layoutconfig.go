package panel

import "git.parallelcoin.io/dev/pod/cmd/gui/vue/mod"

func LayoutConfig() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Layout congifuration",
		ID:       "panellayoutconfig",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "config",
		Js: `
	data () { return { 
	duoSystem,
	restoreModel: []
}},
		methods:{
  onRestore: function(args) {
            // Create instances for DuoDashbord element
             layout.$refs.DuoDashbord.$el.ej2_instances[0].panels = layout.$refs.restoreModel;
        },

        // Save the current panels
        onSave: function(args) {
            // Create instances for DuoDashbord element
            layout.$refs.restoreModel = layout.$refs.DuoDashbord.$el.ej2_instances[0].serialize();
            layout.$refs.restoreModel[0].content = '<div class="content">0</div>';
            layout.$refs.restoreModel[1].content = '<div class="content">1</div>';
            layout.$refs.restoreModel[2].content = '<div class="content">2</div>';
            layout.$refs.restoreModel[3].content = '<div class="content">3</div>';
            layout.$refs.restoreModel[4].content = '<div class="content">4</div>';
            layout.$refs.restoreModel[5].content = '<div class="content">5</div>';
            layout.$refs.restoreModel[6].content = '<div class="content">6</div>';
        },
    btnClick: function(event) {
      if (this.$refs.toggleBtn.$el.classList.contains('e-active')) {
        this.$refs.toggleBtn.content = 'Dark theme';
        this.$refs.toggleBtn.iconCss = 'e-btn-sb-icon e-play-icon';
		document.getElementById('main').classList.toggle('lightTheme');
      } else {
        this.$refs.toggleBtn.content = 'Light theme';
        this.$refs.toggleBtn.iconCss = 'e-btn-sb-icon e-pause-icon';
		document.getElementById('main').classList.toggle('lightTheme');
      }
    },
 		 toggleClick: function(args) {
              if (this.$refs.toggleBtn.$el.textContent == 'Edit') { 
                    layout.$refs.DuoDashbord.allowResizing = true;
                    layout.$refs.DuoDashbord.allowDragging = true;
                    this.$refs.toggleBtn.$el.textContent = 'Save';
                    this.$refs.toggleBtn.iconCss = "save";
                    document.getElementById('dialogBtn').style.layout = 'block';
            } else {
                layout.$refs.DuoDashbord.allowResizing = false;
                layout.$refs.DuoDashbord.allowDragging = false;
                this.$refs.toggleBtn.$el.textContent = 'Edit';
                this.$refs.toggleBtn.iconCss = "edit";
                document.getElementById('dialogBtn').style.layout = 'none';
            	}
			},
        dialogButtonClick: function() {
              layout.$refs.dialogObj.show();
              var proxy = this;
              this.$refs.dialogObj.$el.querySelector('#linetemplate').onclick = function (e) {
                   var panel = {
                       sizeX: 1,
                       sizeY: 1,
                       header: '<div>Line Chart</div>',
                       row: 0,
                       col:0,
                       content: proxy.line
                   }
                   proxy.$refs.DuoDashbord.addPanel(panel);
                   proxy.$refs.dialogObj.hide();
               }
               layout.$refs.dialogObj.$el.querySelector('#pietemplate').onclick = function (e) {
                   var panel = {
                       sizeX: 1,
                       sizeY: 1,
                       header: '<div>Pie Chart</div>',
                       row: 0,
                       col:0,
                       content: proxy.pie
                   }
                   proxy.$refs.DuoDashbord.addPanel(panel);
                   proxy.$refs.dialogObj.hide();
               }
               layout.$refs.dialogObj.$el.querySelector('#splinetemplate').onclick = function (e) {
                   var panel = {
                       sizeX: 2,
                       sizeY: 1,
                       header: '<div>Spline Chart</div>',
                       row: 0,
                       col:0,
                       content: proxy.spline
                   }
                   proxy.$refs.DuoDashbord.addPanel(panel);
                   proxy.$refs.dialogObj.hide();
               }
        } 
        },
		`,
		Template: `<div class="flx flc fii justifyBetween ">
			<ejs-button cssClass="e-outline e-primary flx fii" id="toggleBtn" ref="toggleBtn" iconCss='edit ' isToggle=true v-on:click.native='toggleClick'>Edit</ejs-button>
			<ejs-button class="add-widget-button e-control e-btn e-comp flx fii marginTop" id="dialogBtn" style="display:none" v-on:click="dialogButtonClick($event)">Add</ejs-button>
                    <table id ="remove">
                            <tbody>
                                <tr><td> Properties Panel </td></tr>
                                <tr>
                                    <td>
                                        <!--  Button element declaration -->
                                        <ejs-button id="save" cssClass="e-primary"  v-on:click.native="onSave" >Save Panel</ejs-button>
                                    </td>
                                    <td>
                                        <!--  Button element declaration -->
                                        <ejs-button id="restore" cssClass="e-flat e-outline" v-on:click.native="onRestore">Restore Panel</ejs-button>
                                    </td>
                                </tr>
                            </tbody>
                        </table>


<span v-html="restoreModel"></span>


    <ejs-button ref="toggleBtn" cssClass='e-flat flx fii marginTop' iconCss='e-btn-sb-icon e-play-icon' isToggle=true v-on:click.native='btnClick'>Dark theme</ejs-button>

		</div>`,
		Css: `
		`,
	}
}

//
//func layoutDetailConfig() mod.DuoVUEcomp {
//	return mod.DuoVUEcomp{
//		IsApp:    true,
//		Name:     "Layout congifuration",
//		ID:       "layoutconfig",
//		Version:  "0.0.1",
//		CompType: "panel",
//		SubType:  "config",
//		Js: `
//		data (){
//			return {
//				floatLabelType: 'Never',
//				min: 1,
//				max: 5,
//				rowmin: 0,
//				rowmax: 5,
//				colmin: 0,
//				colmax: 4,
//				value:'',
//				data: ['Panel0', 'Panel1', 'Panel2', 'Panel3', 'Panel4', 'Panel5', 'Panel6'],
//				count: 7
//			};
//		},
//			created: function(args) {
//			  this.$refs.toggleSwitch.toggle();
//			},
//		methods: {
//			onAdd: function(args) {
//				var sizeX = document.getElementById("sizex").ej2_instances["id"];
//				var sizeY = document.getElementById("sizey").ej2_instances["id"];
//				var row = document.getElementById("row").ej2_instances["id"];
//				var column = document.getElementById("column").ej2_instances["id"];
//				var dropdownObject = document.getElementById("dropdown").ej2_instances["id"];
//				var dashboardObj = document.getElementById("dashboard_default").ej2_instances["id"];
//				var panel = [{
//					'id': "Panel"+ this.count.toString(),
//					'sizeX': sizeX.value,
//					'sizeY': sizeY.value,
//					'row': row.value,
//					'col': column.value,
//					'content': "<div class='content'>"+ this.count +"</div>"
//				}];
//				dropdownObject.dataSource.push("Panel" + this.count.toString());
//				dropdownObject.refresh();
//				this.count = this.count + 1;
//				dashboardObj.addPanel(panel[0]);
//			},
//
//			onRemove: function(args) {
//				var dashboardObj = document.getElementById("dashboard_default").ej2_instances["id"];
//				var dropdownObject = document.getElementById("dropdown").ej2_instances["id"];
//				dashboardObj.removePanel(dropdownObject.value);
//				dropdownObject.dataSource.splice(dropdownObject.dataSource.indexOf(dropdownObject.value), 1);
//				dropdownObject.value = null;
//				dropdownObject.refresh();
//			}}
//		`,
//		Template: `<div class="rwrap property-section dashboard" id="api_property">
//        <div className="row property-panel-content">
//            <div className="card-body">
//                <div v-for="panel in core.data.conf.display.layout.panel" className="form-group row">
//                    <table id ="add">
//                        <tbody>
//                            <tr><td id="property"><h4 v-text="panel.header"></h4></td></tr>
//                            <tr>
//                                <td>SizeX</td>
//                                <td>
//                                    <ejs-numerictextbox :id="panel.id + 'sizex'" placeholder="Ex: 1" v-model="panel.sizex" :min="min" :max="max" :floatLabelType="floatLabelType"></ejs-numerictextbox>
//                                </td>
//                            </tr>
//                            <tr>
//                                <td>SizeX</td>
//                                <td>
//                                    <ejs-numerictextbox :id="panel.id + 'sizey'" placeholder="Ex: 1" v-model="panel.sizey" :min="min" :max="max" :floatLabelType="floatLabelType"></ejs-numerictextbox>
//                                </td>
//                            </tr>
//                            <tr>
//                                <td>Row</td>
//                                <td>
//                                    <ejs-numerictextbox :id="panel.id + 'row'" placeholder="Ex: 1" v-model="panel.row" :min="rowmin" :max="rowmax" :floatLabelType="floatLabelType"></ejs-numerictextbox>
//                                </td>
//                            </tr>
//                            <tr>
//                                <td>Column</td>
//                                <td>
//                                    <ejs-numerictextbox :id="panel.id + 'column'" placeholder="Ex: 1" v-model="panel.column" :min="colmin" :max="colmax" :floatLabelType="floatLabelType"></ejs-numerictextbox>
//                                </td>
//                            </tr>
//							<tr>
//                                <td>Enabled</td>
//                                <td>
//                                      <ejs-switch ref="toggleSwitch" checked=true :created="created"></ejs-switch>
//                                </td>
//                            </tr>
//                        </tbody>
//                    </table>
//                    </div>
//                </div>
//            </div>
//		</div>`,
//		Css: `
//
//		`,
//	}
//}
