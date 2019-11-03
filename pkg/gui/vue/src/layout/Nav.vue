<template>
  <nav class="Nav">
    <ul id="menu" class="lsn noPadding">
      <li id='menuoverview' class='sidebar-item current'>
          <button @click="duoSystem.isScreen = 'Overview'" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
            <IcoOverview />
          </button>
      </li>
      <li id='menutransactions' class='sidebar-item'>  
        <button @click="duoSystem.isScreen = 'Transactions'" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
          <IcoHistory />
        </button>
      </li>
      <li id='menuaddressbook' class='sidebar-item'>
        <button @click="duoSystem.isScreen = 'AddressBook'" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
          <IcoAddressBook />
        </button>
      </li>
      <li id='menublockexplorer' class='sidebar-item'>
        <button @click="duoSystem.isScreen = 'Explorer'" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
          <IcoOverview />
        </button>
      </li>
      <li id='menusettings' class='sidebar-item'>
        <button @click="duoSystem.isScreen = 'Settings'" class="noMargin noPadding noBorder bgTrans sXs cursorPointer">
          <IcoSettings />
        </button>
      </li>
    </ul>
  </nav>
</template>

<script>
import { duoSystem } from '../DuOSyS.js'
import IcoOverview from '../ico/IcoOverview.vue'
import IcoHistory from '../ico/IcoHistory.vue'
import IcoAddressBook from '../ico/IcoAddressBook.vue'
import IcoSettings from '../ico/IcoSettings.vue'

export default {
  name: 'Nav',
   data() {
    return {
    duoSystem,
	}},
    components: {
    IcoOverview,
    IcoHistory,
    IcoAddressBook,
    IcoSettings,
  }
}
</script>