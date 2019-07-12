package notifyonce

import (
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestNotifyOnce(t *testing.T) {
	startTime := time.Now()
	deadline := time.NewTicker(Timeout * 2)
	event := fsnotify.Event{Name: "test"}

	// simulate an event that only appears once, it should be communicated via Once
	// after Delay has ticked
	HandleEvent(event)
	select {
	case e := <-Event:
		if time.Now().Sub(startTime) < Timeout {
			t.Error("single notice was sent before Delay time elapsed")
		}
		if e.Name != event.Name {
			t.Errorf("received Event.Name should be `%s`, got `%s`", event.Name, e.Name)
		}
	case <-deadline.C:
		t.Error("single notice was never sent")
	}
	deadline.Stop()

	if _, ok := cancelMap.Get(event.Name); ok {
		t.Errorf("cancelMap should not contain key `%s`", event.Name)
	}

	// simulate an event that occurs twice, it should be communicated once
	startTime = time.Now()
	deadline = time.NewTicker(Timeout)

	HandleEvent(event)
	time.Sleep(Timeout / 2)
	HandleEvent(event)

	select {
	case e := <-Event:
		if time.Now().Sub(startTime) > Timeout {
			t.Error("single notice was sent after Delay time elapsed")
		}
		if e.Name != event.Name {
			t.Errorf("received Event.Name should be `%s`, got `%s`", event.Name, e.Name)
		}
	case <-deadline.C:
		t.Error("single notice was never sent")
	}
	deadline.Stop()

	if _, ok := cancelMap.Get(event.Name); ok {
		t.Errorf("cancelMap should not contain key `%s`", event.Name)
	}
}
