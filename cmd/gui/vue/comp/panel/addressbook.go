package panel

import "github.com/p9c/pod/cmd/gui/vue/mod"

func AddressBook() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Address Book",
		ID:       "paneladdressbook",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "addressbook",
		Js: `
	data () { return { 
	duoSystem,
	account:"default",
	label: "no label",
    pageSettings: { 
		pageSize: 10 
		},
    sortOptions: { 
		columns: [{ 
			field: 'num', direction: 'Ascending' }]
		},
    toolbar: ['Add', 'Edit'],
	labelrules: { 
		required: true 
		},
    editparams: { 
			params: { 
				popupHeight: '300px' 
			}
		},
	editSettings: { 
		allowEditing: true, 
		allowAdding: true, 
		allowDeleting: true, 
		mode: 'Dialog',
		template: function () {
			return { template : partaddress}
      		}
  		}
}},
created: function(){
	
},
methods: {
        actionBegin: function(args) { 
            if(args.requestType == "beginEdit") { 
                this.goCreateAddress(); 
            }; 
            if(args.requestType == "save") { 
                this.goCreateAddress(); 
            } 
        },
		goCreateAddress: function(){
			const addrCmd = {
			account: this.account,
			label: this.label,
			};
			const addrCmdStr = JSON.stringify(addrCmd);
			external.invoke('createAddress:'+addrCmdStr);
		},
},
`,
		Template: `<div class="rwrap">
        <ejs-grid
ref='grid'
height="100%" 
:dataSource='this.duoSystem.addressBook.addresses'
:allowSorting='true' 
:allowPaging='true'
:sortSettings='sortOptions' 
:pageSettings='pageSettings' 
:editSettings='editSettings'
:actionBegin='actionBegin'
:toolbar='toolbar'>
          <e-columns>
            <e-column field='num' headerText='Index' width='80' textAlign='Right' :allowAdding='false' :allowEditing='false'></e-column>
            <e-column field='label' headerText='Label' editType='textedit' :validationRules='labelrules' defaultValue='label' :edit='editparams' textAlign='Right' width=160></e-column>
			<e-column field='address' headerText='Address' textAlign='Right' width=240 :isPrimaryKey='true' :allowEditing='false' :allowAdding='false'></e-column>
            <e-column field='amount' headerText='Amount' textAlign='Right' :allowEditing='false' :allowAdding='false' width=60></e-column>
          </e-columns>
        </ejs-grid>
</div>`,
		Css: `
		`,
	}
}
