package sys

import "git.parallelcoin.io/dev/pod/cmd/gui/vue/mod"

func Dev() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Dev",
		ID:       "dev",
		Version:  "0.0.1",
		CompType: "core",
		SubType:  "dev",
		Js: `
		data () { return { 
			name:   '',
			folder: '',
			data:   '',
			model:  { apps:  core.data.conf.display.layout.items},
			layout: ` + "`" + `{
				fields: [{
					type: "checklist",
					label: "apps",
					model: "apps",
					multi: true,
					required: true,
					multiSelect: true,
					values: core.data.cfg.apps.layout,
					checklistOptions: {
						key: "name",
							id: "i",
							trackBy: "i",
						},
						selectOptions: {
							multiple: true,
							clearOnSelect: false,
							closeOnSelect: false,
						},
					},	{
				   type: "switch",
				   label: "Status",
					model: "status",
					multi: true,
					readonly: false,
					featured: false,
					disabled: false,
					default: true,
					textOn: "Active",
					textOff: "Inactive"
					}]
					}` + "`" + `,
				layoutItem: ` + "`" + `{
					fields: [{
					type: "input",
					inputType: "Number",
					label: "name",
					model: "x",
					min: 1,
					max: 10,
					required: true,
				}]
			}` + "`" + `,
				"formOptions": { validateAfterLoad: true, validateAfterChanged: true },
			}}`,
		Template: `<div class="dev" v-show="duoSystem.isDev">
		<div class="rwrap">
		<div class="dev"><h1>Layout</h1>
		<ul class="flx fcl">
		

		</ul>
		<ul class="flx fcl">
		<li><span v-html="model.apps"></span></li>
		<li class="fii">
		<VueFormGenerator :schema="layout" :model="model" :options="formOptions"></VueFormGenerator>
		</li>

			<li class="fii">
				<ul class="flx fcl">
					<li v-for="app in model.apps" class="fii">
					<VueFormGenerator :schema="layoutItem" :model="app" :options="formOptions"></VueFormGenerator>
					<span v-html="app"></span>
					</li>
				</ul>
			</li>
		</ul>
		</div></div>
		</div>`,
		Css: `
		.dev{
			background:blue;
		}
		`,
	}
}
