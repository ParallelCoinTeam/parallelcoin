<template>
    <main class="pageTransaction">
        <PanelTransactions :transactions="duoSystem.transactions"/>
    </main>
</template>

<script>
import { duoSystem } from '../DuOSyS.js'
import PanelTransactions from '../panels/PanelTransactions.vue'

export default {
  data () { return { 
	duoSystem }},
  name: 'Transactions',
  components: {
    PanelTransactions
  }
}
</script>