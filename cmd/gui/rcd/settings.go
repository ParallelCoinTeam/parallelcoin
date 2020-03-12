package rcd

import (
	"fmt"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/pkg/conte"
	log "github.com/p9c/logi"
	"github.com/p9c/pod/pkg/pod"
	"reflect"
	"time"
)

func (r *RcVar) SaveDaemonCfg() {

	save.Pod(getField(r.cx.Config, r.cx.ConfigMap))

}

func getField(v *pod.Config, configMap map[string]interface{}) *pod.Config {
	//s := reflect.ValueOf(v).Elem()

	for label, data := range configMap {
		s := reflect.ValueOf(v).Elem().FieldByName(label)

		//typeOfT := s.Type()
		//for i := 0; i < s.NumField(); i++ {
		//	f := s.Field(i)

		//fmt.Printf("%d: %s %s = %v\n", i,
		//	typeOfT.Field(i).Name, f.Type(), f.Interface())
		//log.L.Info("lastaviac", f.Type().String())
		if s.IsValid() {
			switch s.Type().String() {
			case "*bool":
				log.L.Info("bool", label)
				log.L.Info("bool", *data.(*bool))
				s.SetBool(*data.(*bool))
				//reflect.ValueOf(&v).Elem().FieldByName(field.Model).SetBool(configMap[field.Model].(bool))
			case "*int":
				log.L.Info("int", label)
				log.L.Info("int", *data.(*int))
				s.SetInt(*data.(*int64))
				//reflect.ValueOf(&v).Elem().FieldByName(field.Model).SetInt(configMap[field.Model].(int64))
			case "*float64":
				log.L.Info("float64", label)
				log.L.Info("float64", *data.(*float64))
				s.SetFloat(*data.(*float64))
				//reflect.ValueOf(&v).Elem().FieldByName(field.Model).SetFloat(configMap[field.Model].(float64))
			case "*string":
				log.L.Info("string", label)
				log.L.Info("string", *data.(*string))
				//s.SetString(*data.(*string))
				//reflect.ValueOf(&v).Elem().FieldByName(field.Model).SetString(configMap[field.Model].(string))
			case "*cli.StringSlice":
				//f.CallSlice(configMap[typeOfT.Field(i).Name].(cli.StringSlice))
				//reflect.ValueOf(&v).Elem().FieldByName(field.Model).Set(configMap[field.Model].(cli.StringSlice))
			case "*time":
				//f.SetBool(*configMap[typeOfT.Field(i).Name].(*bool))
			}
		}
		//log.L.Info("IDE", configMap[typeOfT.Field(i).Name])
		//log.L.Info("IDE", typeOfT.Field(i).Name)
	}
	return v
}

func settings(cx *conte.Xt) *model.DuoUIsettings {

	settings := &model.DuoUIsettings{
		Abbrevation: "DUO",
		Tabs: &model.DuoUIconfTabs{
			Current:  "wallet",
			TabsList: make(map[string]*gel.Button),
		},
		Daemon: &model.DaemonConfig{
			Config: cx.Config,
			Schema: pod.GetConfigSchema(cx.Config, cx.ConfigMap),
		},
	}
	// Settings tabs

	settingsFields := make(map[string]interface{})
	for _, group := range settings.Daemon.Schema.Groups {
		settings.Tabs.TabsList[group.Legend] = new(gel.Button)
		for _, field := range group.Fields {
			switch field.Type {
			case "array":
				settingsFields[field.Label] = new(gel.Button)
			case "input":
				settingsFields[field.Label] = &gel.Editor{
					SingleLine: true,
				}
				if cx.ConfigMap[field.Model] != nil {
					switch field.InputType {
					case "text":
						(settingsFields[field.Label]).(*gel.Editor).SetText(fmt.Sprint(*cx.ConfigMap[field.Model].(*string)))
					case "number":
						(settingsFields[field.Label]).(*gel.Editor).SetText(fmt.Sprint(*cx.ConfigMap[field.Model].(*int)))
					case "decimal":
						(settingsFields[field.Label]).(*gel.Editor).SetText(fmt.Sprint(*cx.ConfigMap[field.Model].(*float64)))
					case "time":
						(settingsFields[field.Label]).(*gel.Editor).SetText(fmt.Sprint(*cx.ConfigMap[field.Model].(*time.Duration)))
					}
				}
			case "switch":
				settingsFields[field.Label] = new(gel.CheckBox)
				(settingsFields[field.Label]).(*gel.CheckBox).SetChecked(*cx.ConfigMap[field.Model].(*bool))
			case "radio":
				settingsFields[field.Label] = new(gel.Enum)
			default:
				settingsFields[field.Label] = new(gel.Button)
			}
		}
	}
	settings.Daemon.Widgets = settingsFields
	return settings
}
