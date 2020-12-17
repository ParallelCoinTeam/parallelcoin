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
	Debug("open channels:", len(createdList), len(createdChannels), occ)
	return o
}

func Ts(n int) C {
	occ := GetOpenChanCount()
	// mx.Lock()
	// defer mx.Unlock()
	createdList = append(createdList, logi.Caller("created quit channel at", 1))
	o := make(C, n)
	createdChannels = append(createdChannels, o)
	Debug("open channels:", occ)
	return o
}

func (c C) Q() {
	loc := GetLocForChan(c)
	// mx.Lock()
	// defer mx.Unlock()
	if !testChanIsClosed(c) {
		close(c)
		Debug("closed channel from " + loc, logi.Caller("from", 1))
	} else {
		Debug("#### channel", loc, "was already closed")
	}
	PrintChanState()
}

func (c C) Wait() <-chan struct{} {
	Debug(logi.Caller(">>>> waiting on quit channel at", 1))
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
			// Debug("still open", createdList[i])
			// createdChannels[i].Q()
		} else {
			Debug(">>>> ", createdList[i])
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
