package sys

import "git.parallelcoin.io/dev/pod/cmd/gui/vue/mod"

func Display() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Display",
		ID:       "display",
		Version:  "0.0.0",
		CompType: "core",
		SubType:  "display",
		Js: `
	  data() {
    return {
		duoSystem,
 		enableDock: true,
        type: 'Over',
        dockSize: '60px',
        aspectRatio: 100/85,
        cellSpacing: [10,10], 
        closeOnDocumentClick: true,
        target: '#sidebarTarget',
	}},
	methods: {
    incr: function() { counter.add(1); },
    blockheightincr: function() { blockheight.addBlockHeight(); },
    btnClick: function(event) {
      if (this.$refs.toggleBtn.$el.classList.contains("e-active")) {
        this.$refs.toggleBtn.content = "FullScreen";
        this.$refs.toggleBtn.iconCss = "e-btn-sb-icons e-play-icon";
		external.invoke('unfullscreen');
      } else {
        this.$refs.toggleBtn.content = "UnFullScreen";
        this.$refs.toggleBtn.iconCss = "e-btn-sb-icons e-pause-icon";
		external.invoke('fullscreen');
      }
    },
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
			setScreen: function (v, m) {
				this.duoSystem.isScreen =	v;
				document.getElementsByClassName('current')[0].classList.remove('current');
                document.getElementById(m).classList.add('current');
				external.invoke('changeTitle:'+ v );
				}
			}`,
		Template: `
    <div id="container">
        <div class="control-section flx fii flc">
    <div class="col-lg-12 col-sm-12 col-md-12 flx flc fii" id="sidebar-section">
      <div id="head">
        <div class="dashboard-header flx justifyBetween itemsCenter">
          <div class="logoContainer flx"><span v-html="system.data.ico.logo" class="flx fii logo"></span>
          </div>
          <div class="searchContent">
            <div class="analysis">ParallelCoin</div>
          </div>
          <div class="right-content">



              <ejs-button
                ref="toggleBtn"
                iconCss="e-btn-sb-icons e-play-icon"
                cssClass="e-small e-flat e-success"
                :isPrimary="true"
                :isToggle="true"
                v-on:click.native="btnClick"
              >FullScreen</ejs-button>

            <div class="information">
              <span id="header-avatar" class="e-avatar e-avatar-medium e-avatar-circle image"></span>
              <div class="text-content">John</div>
            </div>
            
          </div>
        </div>
      </div>
      <!-- sidebar element declaration -->
      <ejs-sidebar id='dashboardSidebar' ref='sidebarInstance' :type='type' :enableDock='enableDock' :dockSize='dockSize' :target='target' :closeOnDocumentClick='false'>
        <div class="content-area">
          <div class="dock">
            <ul>

              <li id='menuoverview' class='sidebar-item current' @click='setScreen("overview","menuoverview");'>
				<span class='e-icons home'></span>
				</li>
              <li id='menutransactions' class='sidebar-item' @click='setScreen("transactions","menutransactions");'>
                <span class='e-icons filter'></span>
              </li>
			<li id='menuaddressbook' class='sidebar-item' @click='setScreen("addressbook","menuaddressbook");'>
                <span class='e-icons filter'></span>
              </li>
              <li id='menublockexplorer' class='sidebar-item' @click='setScreen("blockexplorer","menublockexplorer");'>
                <span class="e-icons analyticsChart"></span>
              </li>
              <li id='menusettings' class='sidebar-item' @click='setScreen("settings","menusettings");'>
					<span class="e-icons settings"></span>
				</li>
              <li id='menucharts' class='sidebar-item' @click='setScreen("charts","menucharts");'>
                <span class="e-icons analytics"></span>
              </li>
            </ul>
          </div>
        </div>
      </ejs-sidebar>
      <!-- end of sidebar element -->
      <!-- main content declaration -->
      <div id="sidebarTarget" class=" flx fii">
        <div class="sidebar-content">
          <div class="dashboardParent">

    <ejs-dashboardlayout 
			v-for="(screen, key) in system.data.conf.display.screens" 
			v-show="duoSystem.isScreen === key" 
			:ref="'DashbordInstance' + key" 
			:columns="9" 
			:id="'Layout' + key" 
			:allowResizing="false"
			:allowDragging="false"
			:cellSpacing="cellSpacing">
			<e-panels>
				<e-panel 
					v-for="panel in screen.panels" 
					:sizeX="panel.sizeX" 
					:sizeY="panel.sizeY" 
					:row="panel.row" 
					:col="panel.col" 
					:header="panel.header"
					:cssClass="panel.cssClass" 
					:content="panel.content" 
					class="rwrap"">
				</e-panel>
            </e-panels>

    </ejs-dashboardlayout>
          </div>
        </div>
      </div>


      <!--end of main content declaration -->
    </div>
  </div>
    </div>
		`,
	}
}
