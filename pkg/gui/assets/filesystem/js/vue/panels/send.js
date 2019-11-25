const PanelSend = {
	data () { return { 
        rcvar }},
    name: 'PanelSend',
    data (){
        return {
            sendAddress:'',
            sendAmount:0.00000000,
            target: '#x',
            header: 'Enter your passphrase',
            content: "<input type='password' value='' class='e-outline fii bgfff' id='passPhrase' required />",
            showCloseIcon: true,
            visible: false,
            buttons: [{
                buttonModel: { isPrimary: true, content: 'Submit' }, click: this.createModalDialog }
            ],
            width: '400px',
            animationSettings: { effect: 'Zoom' },
            animateSettings: { effect: 'Zoom' },
            model_header: 'Verify send details',
            isModal: true,
            mbuttons: [{
                buttonModel: {isPrimary: true, content: 'Yes'}, click:  function() {
                    this.hide();
                    this.goDuoSend();
                }}, {
                buttonModel: {isPrimary: false, content: 'No'}, click: this.showDialog
                }
            ],
            arr: []
		}
	},
	methods: {
        btnClick: function () {
            this.$refs.Dialog.show();
        },
        createModalDialog: function() {
            this.$refs.Dialog.hide();
            this.$refs.ModalDialog.content = this.getDynamicContent();
            this.$refs.ModalDialog.show();
        },
        showDialog: function() {
            this.$refs.ModalDialog.hide();
            this.$refs.Dialog.show();
        },
        getDynamicContent: function() {
            if (!document.getElementById('dialog')) {return;}
            let template = "<div class='row'><div class='col-xs-12 col-sm-12 col-lg-12 col-md-12'><b>Confirm your details</b></div>" +
            "</div><div class='row'><div class='col-xs-6 col-sm-6 col-lg-6 col-md-6'><span id='sendAddress'> Send to address: </span>" +
            "</div><div class='col-xs-6 col-sm-6 col-lg-6 col-md-6'><span id='sendAddressValue'>"+ document.getElementById("sendAddress").value + "</span> </div></div>" +
            "</div>"
            return template;
        },
		goDuoSend: function(){
			const sendCmd = {
			wp: document.getElementById('passPhrase').value,
			ad: this.sendAddress,
			am: this.sendAmount,
		};
		const sendCmdStr = JSON.stringify(sendCmd);

			external.invoke('send:'+sendCmdStr);
		}
	},
		template: `<div class="flx flc fii justifyBetween">


		<ejs-textbox floatLabelType='Auto' cssClass='e-outline flx noMargin fii bgfff' v-model='sendAddress' id='sendAddress' placeholder='Enter DUO address'></ejs-textbox>

<div class="flx fii marginTop">	
<ejs-numerictextbox cssClass="e-outline noMargin fii bgfff" v-model.number="sendAmount" format="c5" validateDecimalOnType=true decimals="8"></ejs-numerictextbox>

<ejs-button id="dialogbtn" v-on:click.native="btnClick">Send</ejs-button>
<ejs-dialog id="dialog" ref="Dialog" :header='header' :target='target' :width='width' :buttons='buttons' :visible='visible' :content='content' :animationSettings='animateSettings'  :showCloseIcon="showCloseIcon">
</ejs-dialog>

<ejs-dialog id="modalDialog" ref="ModalDialog" :header='model_header' :isModal='isModal' :target='target' :buttons='mbuttons' :animationSettings='animationSettings' :visible='visible' >
</ejs-dialog>
</div>
</div>`,
}