package pkgs

import (
	"reflect"
	"strings"
)

func Name(dtype interface{}) string {
	name := reflect.TypeOf(dtype).PkgPath()
	name = strings.TrimPrefix(
		name,
		"github.com/parallelcointeam/parallelcoin/",
	)
	return name
}
