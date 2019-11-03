<template>
  <div id="app" class="fullScreen bgDark flx lightTheme">
    <div id="display" class="fii">
      <div class="grid-container rwrap bgDark">
        <div class="flx fii Logo">
          <Logo />
        </div>
        <Header />
        <div class="Sidebar bgLight">
          <div class="Open"></div>
          <Nav />
          <div class="Side"></div>
        </div>
        <div id="main" class="grayGrad Main">
          <component :is="duoSystem.isScreen"></component>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { duoSystem } from '../DuOSyS.js'
import Logo from '../ico/Logo.vue'
import Header from './Header.vue'
import Nav from './Nav.vue'
import Overview from '../pages/Overview.vue'
import Transactions from '../pages/Transactions.vue'
import AddressBook from '../pages/AddressBook.vue'
import Settings from '../pages/Settings.vue'

export default {
  name: 'Xorg',
  data () { return { 
	duoSystem }},
  components: {
    Logo,
    Header,
    Nav,
    Overview,
    Transactions,
    AddressBook,
    Settings
  }
}
</script>