package util

import (
	"net"
	"strconv"
)

func GetActualPort(listener string) uint16 {
	var e error
	var p string
	if _, p, e = net.SplitHostPort(listener); err.Chk(e) {
	}
	var oI uint64
	if oI, e = strconv.ParseUint(p, 10, 16); err.Chk(e) {
	}
	return uint16(oI)
}
