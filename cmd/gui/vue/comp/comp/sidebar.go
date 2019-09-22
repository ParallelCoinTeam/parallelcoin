package comp

import "github.com/parallelcointeam/parallelcoin/cmd/gui/vue/mod"

func Sidebar() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    false,
		Name:     "Sidebar",
		ID:       "duoSidebar",
		Version:  "0.0.0",
		CompType: "core",
		SubType:  "sidebar",
		Js: `
	  data() {
    return {
		duoSystem,
          dataList: [
                  { slug: 'overview', id: '01', messages: '', badge: '', icon: 'sf-icon-Dashboard' },
                  { slug: 'transactions', id: '02', messages: '', badge: '', icon: 'sf-icon-Hardware' },
                  { slug: 'addressbook', id: '03', messages: '', badge: '', icon: 'sf-icon-Software' },
                  { slug: 'settings', id: '04', messages: '', badge: '', icon: 'sf-icon-License' },
                  { slug: 'explorer', id: '05', messages: '', badge: '', icon: 'sf-icon-Request' },
                  { slug: 'about', id: '06', messages: '', badge: '', icon: 'sf-icon-About' }
            ],
            fields: { iconCss: 'icon', tooltip: 'text' },
            template: function () {
            return { template: duoMenu }
            },
            closeOnDocumentClick: true,
			

	}},
methods:{
 onComplete: function (args) {
            var menuId
            for (let i = 0; i < args.data.length; i++) {
                if (args.data[i].slug === this.duoSystem.activeLayout) {
                menuId = args.data[i].id
                }
            }
            this.$refs.sidebarListObj.selectItem({'id': menuId})
        },
    // Listview select event handler
        onSelect: function (args) {
            this.duoSystem.activeLayout = args.slug;
        },
},

`,
		Template: `
<template>
<duoSidebar class="posRel flx flc bgLight lineRight duoSidebar">
	<div class="content-area">
          <div class="dock">
		<ejs-listview ref="sidebarListObj" id="menulist" :dataSource='dataList' :template="template" :closeOnDocumentClick='closeOnDocumentClick' :fields='fields' :select="onSelect" showIcon='true' :actionComplete="onComplete"></ejs-listview>
	</div>
</duoSidebar>
</template>
`,
		Css: `



`}
}
