package rcd

import (
	js "encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/pod"
)

func (r *RcVar) SaveDaemonCfg() {

	marshalled, _ := js.Marshal(r.Settings.Daemon.Config)
	config := pod.Config{}
	if err := js.Unmarshal(marshalled, &config); err != nil {
	}
	// L.Debug(r.cx.Config)
	// L.Debug(config)
	save.Pod(&config)
}

func getField(v *pod.Config, configMap map[string]interface{}) *pod.Config {
	// s := reflect.ValueOf(v).Elem()

	for label, data := range configMap {
		s := reflect.ValueOf(v).Elem().FieldByName(label)

		// typeOfT := s.Type()
		// for i := 0; i < s.NumField(); i++ {
		//	f := s.Field(i)

		// fmt.Printf("%d: %s %s = %v\n", i,
		//	typeOfT.Field(i).Name, f.Type(), f.Interface())
		// L.Info("lastaviac", f.Type().String())
		if s.IsValid() {
			switch s.Type().String() {
			case "*bool":
				L.Info("bool", label)
				L.Info("bool", *data.(*bool))
				s.SetBool(*data.(*bool))
				// reflect.ValueOf(&v).Elem().FieldByName(field.Model).SetBool(configMap[field.Model].(bool))
			case "*int":
				L.Info("int", label)
				L.Info("int", *data.(*int))
				s.SetInt(*data.(*int64))
				// reflect.ValueOf(&v).Elem().FieldByName(field.Model).SetInt(configMap[field.Model].(int64))
			case "*float64":
				L.Info("float64", label)
				L.Info("float64", *data.(*float64))
				s.SetFloat(*data.(*float64))
				// reflect.ValueOf(&v).Elem().FieldByName(field.Model).SetFloat(configMap[field.Model].(float64))
			case "*string":
				L.Info("string", label)
				L.Info("string", *data.(*string))
				// s.SetString(*data.(*string))
				// reflect.ValueOf(&v).Elem().FieldByName(field.Model).SetString(configMap[field.Model].(string))
			case "*cli.StringSlice":
				// f.CallSlice(configMap[typeOfT.Field(i).Name].(cli.StringSlice))
				// reflect.ValueOf(&v).Elem().FieldByName(field.Model).Set(configMap[field.Model].(cli.StringSlice))
			case "*time":
				// f.SetBool(*configMap[typeOfT.Field(i).Name].(*bool))
			}
		}
		// L.Info("IDE", configMap[typeOfT.Field(i).Name])
		// L.Info("IDE", typeOfT.Field(i).Name)
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
			Config: cx.ConfigMap,
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
				settingsFields[field.Model] = new(gel.Button)
			case "input":
				settingsFields[field.Model] = &gel.Editor{
					SingleLine: true,
				}
				if settings.Daemon.Config[field.Model] != nil {
					switch field.InputType {
					case "text":
						(settingsFields[field.Model]).(*gel.Editor).SetText(fmt.Sprint(*settings.Daemon.Config[field.Model].(*string)))
					case "number":
						(settingsFields[field.Model]).(*gel.Editor).SetText(fmt.Sprint(*settings.Daemon.Config[field.Model].(*int)))
					case "decimal":
						(settingsFields[field.Model]).(*gel.Editor).SetText(fmt.Sprint(*settings.Daemon.Config[field.Model].(*float64)))
					case "time":
						(settingsFields[field.Model]).(*gel.Editor).SetText(fmt.Sprint(*settings.Daemon.Config[field.Model].(*time.Duration)))
					}
				}
			case "switch":
				settingsFields[field.Model] = new(gel.CheckBox)
				(settingsFields[field.Model]).(*gel.CheckBox).SetChecked(*settings.Daemon.Config[field.Model].(*bool))
			case "radio":
				settingsFields[field.Model] = new(gel.Enum)
			default:
				settingsFields[field.Model] = new(gel.Button)
			}
		}
	}
	settings.Daemon.Widgets = settingsFields
	return settings
}
