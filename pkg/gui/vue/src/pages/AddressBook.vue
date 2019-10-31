<template>
    <main class="pageTransaction">
        <PanelAddressBook :addressbook="duoSystem.addressBook"/>
    </main>
</template>

<script>
import { duoSystem } from '../DuOSyS.js'
import PanelAddressBook from '../panels/PanelAddressBook.vue'

export default {
  data () { return { 
	duoSystem }},
  name: 'AddressBook',
  components: {
    PanelAddressBook
  }
}
</script>