package pkgs

import (
	"reflect"
	"strings"
)

func Name(dtype interface{}) string {
	name := reflect.TypeOf(dtype).PkgPath()
	name = strings.TrimPrefix(
		name,
		"git.parallelcoin.io/dev/pod/",
	)
	return name
}
