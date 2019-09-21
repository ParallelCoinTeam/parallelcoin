package serv

import "github.com/parallelcointeam/parallelcoin/cmd/gui/vue/mod"

func Alert() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Alert Services",
		ID:       "srvalert",
		Version:  "0.0.1",
		CompType: "core",
		SubType:  "alert",
		Js: `
	data () { return { 
		duoSystem,
	}},
	watch: {
		'duoSystem.alert.time': function(newVal, oldVal) {
			setTimeout(() => {
				      this.$refs.elementRef.show({
					position:{ X: 'Right', Y: 'Bottom' }, 
					timeOut:3000, 
					title: duoSystem.alert.title, 
					content: duoSystem.alert.message,
					cssClass: 'e-toast-'+duoSystem.alert.type,
					icon: 'e-meeting'
					});
				},200);
		} 
	},
  methods: {
       onClick: function(e){
            e.clickToClose = true;
       },
	beforeOpen: function(e){
          var audio = new Audio('https://drive.google.com/uc?export=download&id=1M95VOpto1cQ4FQHzNBaLf0WFQglrtWi7');
          audio.play();
       }
    }
`,
Template: `<div><ejs-toast ref='elementRef' id='element' :beforeOpen='beforeOpen'></ejs-toast></div>`,
	}
}
