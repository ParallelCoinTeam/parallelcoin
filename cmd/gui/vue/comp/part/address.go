package part

import "github.com/p9c/pod/cmd/gui/vue/mod"

func PartAddress() mod.DuoVUEcomp {
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
