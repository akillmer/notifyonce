package notifyonce

import (
	"context"
	"time"

	"github.com/fsnotify/fsnotify"
	cmap "github.com/orcaman/concurrent-map"
)

var (
	cancelMap = cmap.New()
	// Event sends an event only once
	Event = make(chan *fsnotify.Event) // sends an event only once
	// Timeout waits the specified duration for a second event before sending
	// the original event once. You may need to fine tune this value.
	Timeout = time.Second
)

// HandleEvent handles a received fsnotify.Event reference
func HandleEvent(event fsnotify.Event) {
	if v, ok := cancelMap.Get(event.Name); !ok {
		// this is the first time the event has been seen
		ctx, cancel := context.WithCancel(context.Background())
		cancelMap.Set(event.Name, cancel)
		// if the event is not seen again within the specified Delay duration
		// then the event will be sent via Once
		go timedNotice(ctx, event)
	} else if cancel, ok := v.(context.CancelFunc); ok {
		cancel() // triggers timedNotice to send the Event via Once
	}
}

func timedNotice(ctx context.Context, event fsnotify.Event) {
	var timer = time.NewTicker(Timeout)

	select {
	case <-timer.C:
		break
	case <-ctx.Done():
		break
	}

	timer.Stop()
	Event <- &event
	cancelMap.Remove(event.Name)
}
