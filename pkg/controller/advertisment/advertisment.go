package advertisment

import (
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/simplebuffer"
	"github.com/p9c/pod/pkg/simplebuffer/IPs"
	"github.com/p9c/pod/pkg/simplebuffer/Uint16"
)

func Get(cx *conte.Xt) simplebuffer.Serializers {
	return simplebuffer.Serializers{
		IPs.GetListenable(),
		Uint16.GetPort((*cx.Config.Listeners)[0]),
		Uint16.GetPort((*cx.Config.RPCListeners)[0]),
		Uint16.GetPort(*cx.Config.Controller),
	}
}
