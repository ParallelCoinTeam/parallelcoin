package p2padvt

import (
	"net"
	
	"github.com/niubaoshu/gotiny"
	
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/routeable"
	
	"github.com/p9c/pod/app/conte"
)

var Magic = []byte{'a', 'd', 'v', 1}

type Advertisment struct {
	IP                   net.IP
	P2P, RPC, Controller uint16
}

//
// // LoadContainer takes a message byte slice payload and loads it into a container ready to be decoded
// func LoadContainer(b []byte) (out Container) {
// 	out.Data = b
// 	return
// }

// Get returns an advertisment serializer
func Get(cx *conte.Xt) []byte {
	P2P:=        util.GetActualPort((*cx.Config.P2PListeners)[0])
	RPC:=        util.GetActualPort((*cx.Config.RPCListeners)[0])
	Controller:= util.GetActualPort(*cx.Config.Controller)
	adv := Advertisment{
		IP:         routeable.GetListenable(),
		P2P:        P2P,
		RPC:        RPC,
		Controller: Controller,
	}
	// Debugs(adv)
	ad := gotiny.Marshal(&adv)
	// Debugs(ad)
	return ad
	// return simplebuffer.Serializers{
	// 	IPs.GetListenable(),
	// 	Uint16.GetPort((*cx.Config.P2PListeners)[0]),
	// 	Uint16.GetPort((*cx.Config.RPCListeners)[0]),
	// 	Uint16.GetPort(*cx.Config.Controller),
	// }
}

// GetAdvt returns an advertisment serializer
func GetAdvt(cx *conte.Xt) *Advertisment {
	adv := &Advertisment{
		IP:         routeable.GetListenable(),
		P2P:        util.GetActualPort((*cx.Config.P2PListeners)[0]),
		RPC:        util.GetActualPort((*cx.Config.RPCListeners)[0]),
		Controller: util.GetActualPort(*cx.Config.Controller),
	}
	return adv
	// return simplebuffer.Serializers{
	// 	IPs.GetListenable(),
	// 	Uint16.GetPort((*cx.Config.P2PListeners)[0]),
	// 	Uint16.GetPort((*cx.Config.RPCListeners)[0]),
	// 	Uint16.GetPort(*cx.Config.Controller),
	// }
}

//
// // GetIPs decodes the IPs from the advertisment
// func (j *Container) GetIPs() []*net.IP {
// 	return IPs.New().DecodeOne(j.Get(0)).Get()
// }
//
// // GetP2PListenersPort returns the p2p listeners port from the advertisment
// func (j *Container) GetP2PListenersPort() uint16 {
// 	return Uint16.New().DecodeOne(j.Get(1)).Get()
// }
//
// // GetRPCListenersPort returns the RPC listeners port from the advertisment
// func (j *Container) GetRPCListenersPort() uint16 {
// 	return Uint16.New().DecodeOne(j.Get(2)).Get()
// }
//
// // GetControllerListenerPort returns the controller listener port from the
// // advertisment
// func (j *Container) GetControllerListenerPort() uint16 {
// 	return Uint16.New().DecodeOne(j.Get(3)).Get()
// }
