package qu

import (
	"github.com/p9c/pod/pkg/logg"
	"strings"
	"sync"
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
	createdList = append(createdList, logg.Caller("chan from", 1))
	o := make(C)
	createdChannels = append(createdChannels, o)
	// F.Ln("open channels:", len(createdList), len(createdChannels), occ)
	return o
}

func Ts(n int) C {
	// PrintChanState()
	// occ := GetOpenChanCount()
	mx.Lock()
	defer mx.Unlock()
	createdList = append(createdList, logg.Caller("buffered chan at", 1))
	o := make(C, n)
	createdChannels = append(createdChannels, o)
	// F.Ln("open channels:", len(createdList), len(createdChannels), occ)
	return o
}

func (c C) Q() {
	loc := GetLocForChan(c)
	mx.Lock()
	if !testChanIsClosed(c) {
		F.Ln("closing chan from "+loc, logg.Caller("\n"+strings.Repeat(" ", 48)+"from", 1))
		close(c)
	} else {
		F.Ln(
			"from"+logg.Caller("", 1), "\n"+strings.Repeat(" ", 48)+
				"channel", loc, "was already closed",
		)
	}
	mx.Unlock()
	// PrintChanState()
}

func (c C) Signal() {
	c <- struct{}{}
}

func (c C) Wait() <-chan struct{} {
	// F.Ln(logg.Caller(">>> waiting on quit channel at", 1))
	return c
}

func testChanIsClosed(ch C) (o bool) {
	if ch == nil {
		return true
	}
	select {
	case <-ch:
		// D.Ln("chan is closed")
		o = true
	default:
	}
	// D.Ln("chan is not closed")
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
	D.Ln("cleaning up closed channels (more than 50 now closed)")
	var c []C
	var l []string
	// D.Ln(">>>>>>>>>>>")
	for i := range createdChannels {
		if i >= len(createdList) {
			break
		}
		if testChanIsClosed(createdChannels[i]) {
			// F.Ln(">>> closed", createdList[i])
			// createdChannels[i].Q()
		} else {
			c = append(c, createdChannels[i])
			l = append(l, createdList[i])
			// F.Ln("<<< open", createdList[i])
		}
		// D.Ln(">>>>>>>>>>>")
	}
	createdChannels = c
	createdList = l
}

func PrintChanState() {
	D.Ln(">>>>>>>>>>>")
	for i := range createdChannels {
		if i >= len(createdList) {
			break
		}
		if testChanIsClosed(createdChannels[i]) {
			F.Ln(">>> closed", createdList[i])
			// createdChannels[i].Q()
		} else {
			F.Ln("<<< open", createdList[i])
		}
	}
	D.Ln(">>>>>>>>>>>")
}

func GetOpenChanCount() (o int) {
	mx.Lock()
	// D.Ln(">>>>>>>>>>>")
	var c int
	for i := range createdChannels {
		if i >= len(createdChannels) {
			break
		}
		if testChanIsClosed(createdChannels[i]) {
			// D.Ln("still open", createdList[i])
			// createdChannels[i].Q()
			c++
		} else {
			o++
			// D.Ln(">>>> ",createdList[i])
		}
		// D.Ln(">>>>>>>>>>>")
	}
	if c > 50 {
		RemoveClosedChans()
	}
	mx.Unlock()
	// o -= len(createdChannels)
	return
}
