<script>
	import { duoSystem } from '../../firmware.js'
	import IcoOverview from "../ico/IcoOverview.svelte";
	import IcoHistory from "../ico/IcoHistory.svelte";
	import IcoAddressBook from "../ico/IcoAddressBook.svelte";
	import IcoSettings from "../ico/IcoSettings.svelte";


</script>




<nav class="Nav">
    <ul id="menu" class="lsn noPadding">
      <li id='menuoverview' class='sidebar-item current'>
          <button on:click={ duoSystem.isScreen = 'Overview'} class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
            <IcoOverview />
          </button>
      </li>
      <li id='menutransactions' class='sidebar-item'>
        <button on:click={ duoSystem.isScreen = 'Transactions'} class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
          <IcoHistory />
        </button>
      </li>
      <li id='menuaddressbook' class='sidebar-item'>
        <button on:click={ duoSystem.isScreen = 'AddressBook'} class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
          <IcoAddressBook />
        </button>
      </li>
      <li id='menublockexplorer' class='sidebar-item'>
        <button on:click={ duoSystem.isScreen = 'Explorer'} class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
          <IcoOverview />
        </button>
      </li>
      <li id='menusettings' class='sidebar-item'>
        <button on:click={ duoSystem.isScreen = 'Settings'} class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
          <IcoSettings />
        </button>
      </li>
    </ul>
</nav>

<style>
#menu svg{
    width:48px;
    height:48px;
}
</style>