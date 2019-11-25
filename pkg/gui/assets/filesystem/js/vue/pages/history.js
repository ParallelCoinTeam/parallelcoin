var PageHistory = {
	data () { return { 
	  rcvar }},
	name: 'History',
  	template: `<main class="pageTransaction">
	  <PanelHistory :transactions="duoSystem.transactions"/>
  </main>`,
  components: {
    PanelHistory
  }
}