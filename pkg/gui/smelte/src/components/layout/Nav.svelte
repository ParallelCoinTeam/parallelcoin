<script>
	import { isPage } from '../../stores/pages.js'
	import IcoOverview from "../ico/IcoOverview.svelte";
	import IcoHistory from "../ico/IcoHistory.svelte";
	import IcoAddressBook from "../ico/IcoAddressBook.svelte";
	import IcoSettings from "../ico/IcoSettings.svelte";

</script>



<nav class="Nav">
    <ul id="menu" class="lsn noPadding">
      <li id='menuoverview' class='sidebar-item current'>
          <button  on:click={isPage.overview} class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
          <IcoOverview />
          </button>
      </li>
      <li id='menutransactions' class='sidebar-item'>
        <button  on:click={isPage.transactions} class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
          <IcoHistory />
        </button>
      </li>
      <li id='menuaddressbook' class='sidebar-item'>
        <button  on:click={isPage.addressbook} class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
            <IcoAddressBook />
        </button>
      </li>
      <li id='menublockexplorer' class='sidebar-item'>
        <button  on:click={isPage.explorer} class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
        <IcoOverview />
        </button>
      </li>
      <li id='menusettings' class='sidebar-item'>
        <button  on:click={isPage.settings} class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
          <IcoSettings />
        </button>
      </li>
    </ul>
</nav>

<style>

</style>