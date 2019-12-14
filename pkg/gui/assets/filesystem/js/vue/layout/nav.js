Vue.component('Nav', {
	name: 'Nav',
  	template: `<nav class="Nav textCenter justifyCenter">
	  <ul id="menu" class="lsn noPadding">
		<li id='menuoverview' class='sidebar-item current'>
			<button @click="duOSnav.getScreen('PageOverview')" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
			  <IcoOverview />
			</button>
		</li>
		<li id='menutransactions' class='sidebar-item'>  
		  <button @click="duOSnav.getScreen('PageHistory')" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
			<IcoHistory />
		  </button>
		</li>
		<li id='menuaddressbook' class='sidebar-item'>
		  <button @click="duOSnav.getScreen('PageAddressBook')" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
			<IcoAddressBook />
		  </button>
		</li>
		<li id='menublockexplorer' class='sidebar-item'>
		  <button @click="duOSnav.getScreen('PageExplorer')" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
			<IcoExplorer />
		  </button>
		</li>
		<li id='menusettings' class='sidebar-item'>
		  <button @click="duOSnav.getScreen('PageSettings')" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
			<IcoSettings />
		  </button>
		</li>
	  </ul>
	</nav>`,
});