var Nav = {
	name: 'Nav',
	data() {
	 return {
	 rcvar,
	 }},
  	template: `<nav class="Nav textCenter justifyCenter">
	  <ul id="menu" class="lsn noPadding">
		<li id='menuoverview' class='sidebar-item current'>
			<button @click="rcvar.isScreen = 'PageOverview'" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
			  <IcoOverview />
			</button>
		</li>
		<li id='menutransactions' class='sidebar-item'>  
		  <button @click="rcvar.isScreen = 'PageTransactions'" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
			<IcoHistory />
		  </button>
		</li>
		<li id='menuaddressbook' class='sidebar-item'>
		  <button @click="rcvar.isScreen = 'PageAddressBook'" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
			<IcoAddressBook />
		  </button>
		</li>
		<li id='menublockexplorer' class='sidebar-item'>
		  <button @click="rcvar.isScreen = 'PageExplorer'" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
			<IcoExplorer />
		  </button>
		</li>
		<li id='menusettings' class='sidebar-item'>
		  <button @click="rcvar.isScreen = 'PageSettings'" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
			<IcoSettings />
		  </button>
		</li>
	  </ul>
	</nav>`,
  	components: {
		IcoOverview,
		IcoHistory,
		IcoExplorer,
		IcoAddressBook,
		IcoSettings,
	}
}