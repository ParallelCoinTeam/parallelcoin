package panel

import "github.com/p9c/pod/cmd/gui/OLDbkp/vue/mod"

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
			return { template : partaddress}
      		}
  		}
}},
created: function(){
	
},
methods: {
        actionBegin: function(args) { 
        	if (args.requestType === 'add') {
				this.goCreateAddress();
				args.data.address = createAddress;
            }; 
            if(args.requestType == "save") { 
                this.goSaveAddressLabel(); 
            } 
        },
	    actionComplete(args) {
			if (args.requestType === 'add') {
            };
        },
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
`,
		Template: `<div class="rwrap">
        <ejs-grid
ref='grid'
height="100%" 
:dataSource='this.duoSystem.addressbook.addresses'
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
</div>`,
		Css: `
		`,
	}
}

func Address() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    false,
		Name:     "Part Address",
		ID:       "partaddress",
		Version:  "0.0.1",
		CompType: "part",
		SubType:  "dialog",
		Js: `
	data () { return {
		duoSystem,
 		width: "300px",
		height: "300px",
		mode: "SVG",
		displaytext: { visibility: false },
		type: "QRCode",
		fontcolorvalue: "#303030",
	errorCorrectionLevelSrc: 3,
}},
 methods: { 
  } 
`,
		Template: `
  <div formGroup="orderForm">

        <div class="form-row flx justifyCenter">

          <ejs-qrcodegenerator
              id="barcode"
              ref="barcodeControl"
              :width="width"
              :displayText="displaytext"
              :height="height"
              :value="$data.data.address"
              :mode="mode"
            ></ejs-qrcodegenerator>

</div>
        <div class="form-row">
            <div class="form-group col-md-6">
                <div class="e-float-input e-control-wrapper">
                    
                    <span class="e-float-line"></span>
  <h3 v-html="$data.data.address"></h3>
                </div>
            </div>
        </div>

<div class="form-row">
            <div class="form-group col-md-6">
                <div class="e-float-input e-control-wrapper">
                    <input ref='label' id="label" name="label" v-model="$data.data.label" type="text">
                    <span class="e-float-line"></span>
  <span v-html="$data.data.label"></span>

                    <label class="e-float-text e-label-top" for="num"> Label</label>
                </div>
            </div>
        </div>


    </div>
`,
	}
}
