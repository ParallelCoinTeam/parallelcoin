package qu

import (
	"sync"

	"github.com/p9c/pod/pkg/util/logi"
)

type C chan struct{}

var createdList []string
var createdChannels []C

var mx sync.Mutex

func T() C {
	// PrintChanState()
	// occ := GetOpenChanCount()
	mx.Lock()
	defer mx.Unlock()
	createdList = append(createdList, logi.Caller("chan from", 1))
	o := make(C)
	createdChannels = append(createdChannels, o)
	// Trace("open channels:", len(createdList), len(createdChannels), occ)
	return o
}

func Ts(n int) C {
	// PrintChanState()
	// occ := GetOpenChanCount()
	mx.Lock()
	defer mx.Unlock()
	createdList = append(createdList, logi.Caller("buffered chan at", 1))
	o := make(C, n)
	createdChannels = append(createdChannels, o)
	// Trace("open channels:", len(createdList), len(createdChannels), occ)
	return o
}

func (c C) Q() {
	loc := GetLocForChan(c)
	mx.Lock()
	if !testChanIsClosed(c) {
		Trace("closing chan from "+loc, logi.Caller("from", 1))
		close(c)
	} else {
		Trace("#### channel", loc, "was already closed")
	}
	mx.Unlock()
	// PrintChanState()
}

func (c C) Wait() <-chan struct{} {
	Trace(logi.Caller(">>> waiting on quit channel at", 1))
	return c
}

func testChanIsClosed(ch C) (o bool) {
	if ch == nil {
		return true
	}
	select {
	case <-ch:
		// Debug("chan is closed")
		o = true
	default:
	}
	// Debug("chan is not closed")
	return
}

func GetLocForChan(c C) (s string) {
	s = "not found"
	mx.Lock()
	for i := range createdList {
		if i >= len(createdChannels) {
			break
		}
		if createdChannels[i] == c {
			s = createdList[i]
		}
	}
	mx.Unlock()
	return
}


func RemoveClosedChans() {
	Debug("cleaning up closed channels (more than 50 now closed)")
	var c []C
	var l []string
	// Debug(">>>>>>>>>>>")
	for i := range createdChannels {
		if i >= len(createdList) {
			break
		}
		if testChanIsClosed(createdChannels[i]) {
			// Trace(">>> closed", createdList[i])
			// createdChannels[i].Q()
		} else {
			c = append(c, createdChannels[i])
			l = append(l, createdList[i])
			// Trace("<<< open", createdList[i])
		}
		// Debug(">>>>>>>>>>>")
	}
	createdChannels = c
	createdList = l
}

func PrintChanState() {

	// Debug(">>>>>>>>>>>")
	for i := range createdChannels {
		if i >= len(createdList) {
			break
		}
		if testChanIsClosed(createdChannels[i]) {
			Trace(">>> closed", createdList[i])
			// createdChannels[i].Q()
		} else {
			Trace("<<< open", createdList[i])
		}
		// Debug(">>>>>>>>>>>")
	}
}

func GetOpenChanCount() (o int) {
	mx.Lock()
	// Debug(">>>>>>>>>>>")
	var c int
	for i := range createdChannels {
		if i >= len(createdChannels) {
			break
		}
		if testChanIsClosed(createdChannels[i]) {
			// Debug("still open", createdList[i])
			// createdChannels[i].Q()
			c++
		} else {
			o++
			// Debug(">>>> ",createdList[i])
		}
		// Debug(">>>>>>>>>>>")
	}
	if c > 50 {
		RemoveClosedChans()
	}
	mx.Unlock()
	// o -= len(createdChannels)
	return
}
