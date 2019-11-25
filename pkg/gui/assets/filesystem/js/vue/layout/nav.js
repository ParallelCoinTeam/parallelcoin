Vue.component('Nav', {
	name: 'Nav',
	data() {
	 return {
	 rcvar, nav,
	 }},
  	template: `<nav class="Nav textCenter justifyCenter">
	  <ul id="menu" class="lsn noPadding">
		<li id='menuoverview' class='sidebar-item current'>
			<button @click="nav.getScreen('PageOverview')" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
			  <IcoOverview />
			</button>
		</li>
		<li id='menutransactions' class='sidebar-item'>  
		  <button @click="nav.getScreen('PageHistory')" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
			<IcoHistory />
		  </button>
		</li>
		<li id='menuaddressbook' class='sidebar-item'>
		  <button @click="nav.getScreen('PageAddressBook')" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
			<IcoAddressBook />
		  </button>
		</li>
		<li id='menublockexplorer' class='sidebar-item'>
		  <button @click="nav.getScreen('PageExplorer')" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
			<IcoExplorer />
		  </button>
		</li>
		<li id='menusettings' class='sidebar-item'>
		  <button @click="nav.getScreen('PageSettings')" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
			<IcoSettings />
		  </button>
		</li>
	  </ul>
	</nav>`,
});