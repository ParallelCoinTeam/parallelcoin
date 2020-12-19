package qu

import (
	"github.com/p9c/pod/pkg/util/logi"
)

type C chan struct{}

var createdList []string
var createdChannels []C

// var mx sync.Mutex

func T() C {
	PrintChanState()
	occ := GetOpenChanCount()
	// mx.Lock()
	// defer mx.Unlock()
	createdList = append(createdList, logi.Caller("remaining quit channel from", 1))
	o := make(C)
	createdChannels = append(createdChannels, o)
	Trace("open channels:", len(createdList), len(createdChannels), occ)
	return o
}

func Ts(n int) C {
	PrintChanState()
	occ := GetOpenChanCount()
	// mx.Lock()
	// defer mx.Unlock()
	createdList = append(createdList, logi.Caller("remaining buffered quit channel at", 1))
	o := make(C, n)
	createdChannels = append(createdChannels, o)
	Trace("open channels:", len(createdList), len(createdChannels), occ)
	return o
}

func (c C) Q() {
	loc := GetLocForChan(c)
	// mx.Lock()
	// defer mx.Unlock()
	if !testChanIsClosed(c) {
		close(c)
		Trace("closing channel from " + loc, logi.Caller("from", 1))
	} else {
		Trace("#### channel", loc, "was already closed")
	}
	PrintChanState()
}

func (c C) Wait() <-chan struct{} {
	Trace(logi.Caller(">>>> waiting on quit channel at", 1))
	return c
}

func testChanIsClosed(ch C) (o bool) {
	select {
	case <-ch:
		o = true
	default:
	}
	return
}

func GetLocForChan(c C) string {
	// mx.Lock()
	// defer mx.Unlock()
	for i := range createdList {
		if i >= len(createdChannels) {
			break
		}
		if createdChannels[i] == c {
			return createdList[i]
		}
	}
	return "not found"
}

func PrintChanState() {
	// mx.Lock()
	// defer mx.Unlock()
	
	// Debug(">>>>>>>>>>>")
	for i := range createdChannels {
		if i >= len(createdList) {
			break
		}
		if testChanIsClosed(createdChannels[i]) {
			Trace("closed", createdList[i])
			// createdChannels[i].Q()
		} else {
			Trace("open", createdList[i])
		}
		// Debug(">>>>>>>>>>>")
	}
}

func GetOpenChanCount() (o int) {
	// mx.Lock()
	// defer mx.Unlock()
	// Debug(">>>>>>>>>>>")
	for i := range createdChannels {
		if i >= len(createdChannels) {
			break
		}
		if testChanIsClosed(createdChannels[i]) {
			// Debug("still open", createdList[i])
			// createdChannels[i].Q()
		} else {
			o++
			// Debug(">>>> ",createdList[i])
		}
		// Debug(">>>>>>>>>>>")
	}
	// o -= len(createdChannels)
	return
}
