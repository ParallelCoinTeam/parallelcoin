package util

import (
	"net"
	"strconv"
)

func GetActualPort(listener string) uint16 {
	var err error
	var p string
	if _, p, err = net.SplitHostPort(listener); Check(err) {
	}
	var oI uint64
	if oI, err = strconv.ParseUint(p, 10, 16); Check(err) {
	}
	return uint16(oI)
}
