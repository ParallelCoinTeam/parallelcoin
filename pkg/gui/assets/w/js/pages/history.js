var PageHistory = {
	data () { return { 
	  duoSystem }},
	name: 'History',
  	template: `<main class="pageTransaction">
	  <PanelHistory :transactions="duoSystem.transactions"/>
  </main>`,
  components: {
    PanelHistory
  }
}