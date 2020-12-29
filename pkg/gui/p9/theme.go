package p9

import (
	"gioui.org/text"
	"gioui.org/unit"
	qu "github.com/p9c/pod/pkg/util/quit"
)

type CallbackQueue chan func()

func NewCallbackQueue(bufSize int) CallbackQueue {
	return make(CallbackQueue, bufSize)
}

type Theme struct {
	quit                      qu.C
	shaper                    text.Shaper
	collection                []text.FontFace
	TextSize                  unit.Value
	Colors                    Colors
	icons                     map[string]*Icon
	scrollBarSize             int
	Dark                      *bool
	iconCache                 IconCache
	WidgetPool                *Pool
	BackgroundProcessingQueue CallbackQueue
}

// NewTheme creates a new theme to use for rendering a user interface
func NewTheme(fontCollection []text.FontFace, quit qu.C) (th *Theme) {
	th = &Theme{
		quit:          quit,
		shaper:        text.NewCache(fontCollection),
		collection:    fontCollection,
		TextSize:      unit.Sp(16),
		Colors:        NewColors(),
		scrollBarSize: 0,
		iconCache:     make(IconCache),
		// 32 should buffer all pending events without blocking during normal operation
		BackgroundProcessingQueue: NewCallbackQueue(32),
		
	}
	th.WidgetPool = th.NewPool()
	// callback channel handler
	go func() {
		Debug("starting background processing queue")
	out:
		for {
			select {
			case fn := <-th.BackgroundProcessingQueue:
				Debug("running background task")
				// this loop runs everything sequentially, the buffer is set to cope with occasional buildups
				fn()
			case <-th.quit:
				Debug("quitting background processing queue")
				break out
			}
		}
	}()
	return
}
