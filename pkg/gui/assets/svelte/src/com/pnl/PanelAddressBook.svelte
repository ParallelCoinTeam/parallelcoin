<div class="rwrap">
        <grid
ref='grid'
height="100%" 
:dataSource='this.duOSys.addressbook.addresses'
:allowSorting='true' 
:allowPaging='true'
:sortSettings='sortOptions' 
:pageSettings='pageSettings' 
:editSettings='editSettings'
:actionBegin='actionBegin'
:actionComplete='actionComplete'
:toolbar='toolbar'>
          <columns>
            <column field='num' headerText='Index' width='80' textAlign='Right' :allowAdding='false' :allowEditing='false'></column>
            <column field='label' headerText='Label' editType='textedit' :validationRules='labelrules' defaultValue='label' :edit='editparams' textAlign='Right' width=160></column>
			<column field='address' headerText='Address' textAlign='Right' width=240 :isPrimaryKey='true' :allowEditing='false' :allowAdding='false'></column>
            <column field='amount' headerText='Amount' textAlign='Right' :allowEditing='false' :allowAdding='false' width=60></column>
          </columns>
        </grid>
</div>

