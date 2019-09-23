package sys

import "github.com/p9c/pod/cmd/gui/vue/mod"

func Alerts() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Alerts",
		ID:       "alerts",
		Version:  "0.0.1",
		CompType: "core",
		SubType:  "alerts",
		Js: `
	data () { return { 
	duoSystem,
            position: {
                X: 'Right',
				Y: 'Bottom'
            }
        }
    },
	created: function() {
		alert.getAlert();
		this.duoSystem.timer = setInterval(alert.getAlert, 500);

	},
	watch: {
	  'alert.data.time': function(newVal, oldVal) {
		setTimeout(() => {
            this.$refs.defaultRef.show({
                title: 'Adaptive Tiles Meeting', content: 'Conference Room 01 / Building 135 10:00 AM-10:30 AM',
                icon: 'e-meeting'
            	});
        	},200);
		}
	},
    mounted: function() {
        setTimeout(() => {
            this.$refs.defaultRef.show({
                title: 'Adaptive Tiles Meeting', content: 'Conference Room 01 / Building 135 10:00 AM-10:30 AM',
                icon: 'e-meeting'
            });
        },200);
    },
		`,
		Template: `<div><ejs-toast ref='defaultRef' id='toast_default' :created="created" :position='position'></ejs-toast></div>`,
	}
}
