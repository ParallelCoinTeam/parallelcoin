package rcd

import (
	"strings"

	log "github.com/p9c/logi"

	"github.com/p9c/pod/pkg/rpc/btcjson"
)

func (r *RcVar) ConsoleCmd(com string) (o string) {
	split := strings.Split(com, " ")
	params := make([]interface{}, 0, len(split[1:]))
	for _, arg := range split[1:] {
		// if arg == "-" {
		// 	param, err := bio.ReadString('\n')
		// 	if err != nil && err != io.EOF {
		// 		fmt.Fprintf(os.Stderr,
		// 			"Failed to read data from stdin: %v\n", err)
		// 		os.Exit(1)
		// 	}
		// 	if err == io.EOF && len(param) == 0 {
		// 		fmt.Fprintln(os.Stderr, "Not enough lines provided on stdin")
		// 		os.Exit(1)
		// 	}
		// 	param = strings.TrimRight(param, "\r\n")
		// 	params = append(params, param)
		// 	continue
		// }
		params = append(params, arg)
	}
	// if err != nil {
	// 	o = fmt.Sprint(err)
	// }
	// method := btcjson.GetMethods()
	// log.L.Info(split)
	// log.L.Infos(cmd)
	var cmd interface{}
	var err error
	if cmd, err = btcjson.NewCmd(split[0], params...); log.L.Check(err) {
	}
	log.L.Debugs(cmd)
	var b []byte
	if b, err = btcjson.MarshalCmd(nil, cmd); log.L.Check(err) {
	}
	log.L.Debug(string(b))
	return
}
