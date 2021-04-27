package wallet

import (
	"strings"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/btcjson"
	"github.com/p9c/matrjoska/pkg/rpchelp"
)

func serverMethods() map[string]struct{} {
	m := make(map[string]struct{})
	for method, handlerData := range RPCHandlers {
		if !handlerData.NoHelp {
			m[method] = struct{}{}
		}
	}
	return m
}

// TestRPCMethodHelpGeneration ensures that help text can be generated for every method of the RPC server for every
// supported locale.
func TestRPCMethodHelpGeneration(t *testing.T) {
	needsGenerate := false
	defer func() {
		if needsGenerate && !t.Failed() {
			t.Error("Generated help texts are out of date: run 'go generate'")
			return
		}
		if t.Failed() {
			t.Log("Regenerate help texts with 'go generate' after fixing")
		}
	}()
	for i := range rpchelp.HelpDescs {
		svrMethods := serverMethods()
		locale := rpchelp.HelpDescs[i].Locale
		generatedDescs := LocaleHelpDescs[locale]()
		for _, m := range rpchelp.Methods {
			delete(svrMethods, m.Method)
			helpText, e := btcjson.GenerateHelp(m.Method, rpchelp.HelpDescs[i].Descs, m.ResultTypes...)
			if e != nil {
				t.Errorf("Cannot generate '%s' help for method '%s': missing description for '%s'",
					locale, m.Method, e,
				)
				continue
			}
			if !needsGenerate && helpText != generatedDescs[m.Method] {
				needsGenerate = true
			}
		}
		for m := range svrMethods {
			t.Errorf("Missing '%s' help for method '%s'", locale, m)
		}
	}
}

// TestRPCMethodUsageGeneration ensures that single line usage text can be generated for every supported request of the
// RPC server.
func TestRPCMethodUsageGeneration(t *testing.T) {
	needsGenerate := false
	defer func() {
		if needsGenerate && !t.Failed() {
			t.Error("Generated help usages are out of date: run 'go generate'")
			return
		}
		if t.Failed() {
			t.Log("Regenerate help usage with 'go generate' after fixing")
		}
	}()
	svrMethods := serverMethods()
	usageStrs := make([]string, 0, len(rpchelp.Methods))
	for _, m := range rpchelp.Methods {
		delete(svrMethods, m.Method)
		usage, e := btcjson.MethodUsageText(m.Method)
		if e != nil {
			t.Errorf("Cannot generate single line usage for method '%s': %v",
				m.Method, e,
			)
		}
		if !t.Failed() {
			usageStrs = append(usageStrs, usage)
		}
	}
	if !t.Failed() {
		usages := strings.Join(usageStrs, "\n")
		needsGenerate = usages != RequestUsages
	}
}
