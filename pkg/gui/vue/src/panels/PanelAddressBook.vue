<script>
import { duoSystem } from '../DuOSyS.js'
import PartAddress from '../parts/PartAddress.vue'

export default {
	data () { return { 
	duoSystem,
	address:"",
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
			return { PartAddress}
			}
		}
}},
created: function(){
	
},
methods: {
        actionBegin: function(args) { 
			if (args.requestType === 'add') {
				this.goCreateAddress();
				// args.data.address = createAddress;
			}
            if(args.requestType == "save") { 
                this.goSaveAddressLabel(); 
            } 
        },
		// actionComplete(args) {
		// 	if (args.requestType === 'add') {
        //     }
        // },
		goCreateAddress: function(){
			const addrCmd = {
			account: this.account,
			};
			const addrCmdStr = JSON.stringify(addrCmd);
			external.invoke('createAddress:'+addrCmdStr);
		},
		goSaveAddressLabel: function(){
			const addrCmd = {
			address: this.address,
			label: this.label,
			};
			const addrCmdStr = JSON.stringify(addrCmd);
			external.invoke('saveAddressLabel:'+addrCmdStr);
		},
},
}
</script>
<template>
<div class="rwrap">
        <ejs-grid
ref='grid'
height="100%" 
:dataSource='this.addressbook.addresses'
:allowSorting='true' 
:allowPaging='true'
:sortSettings='sortOptions' 
:pageSettings='pageSettings' 
:editSettings='editSettings'
:actionBegin='actionBegin'
:actionComplete='actionComplete'
:toolbar='toolbar'>
          <e-columns>
            <e-column field='num' headerText='Index' width='80' textAlign='Right' :allowAdding='false' :allowEditing='false'></e-column>
            <e-column field='label' headerText='Label' editType='textedit' :validationRules='labelrules' defaultValue='label' :edit='editparams' textAlign='Right' width=160></e-column>
			<e-column field='address' headerText='Address' textAlign='Right' width=240 :isPrimaryKey='true' :allowEditing='false' :allowAdding='false'></e-column>
            <e-column field='amount' headerText='Amount' textAlign='Right' :allowEditing='false' :allowAdding='false' width=60></e-column>
          </e-columns>
        </ejs-grid>
</div>
</template>