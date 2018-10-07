package watch

import (
	"sync"

	"github.com/fvdveen/mu2-config/events"
	searchpb "github.com/fvdveen/mu2-proto/go/proto/search"
	"github.com/fvdveen/mu2-search/search"
	"github.com/sirupsen/logrus"
)

// Youtube creates a watcher for youtube events
func Youtube(ch <-chan *events.Event) searchpb.SearchServiceHandler {
	var wg sync.WaitGroup

	s := search.NewWrapper(nil)

	wg.Add(1)

	go func(ch <-chan *events.Event, s search.Wrapper) {
		logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "youtube"}).Debug("Starting...")
		var done = false
		for evnt := range ch {
			if evnt.Key != "youtube.apikey" {
				continue
			}

			srv, err := search.New(evnt.Change)
			if err != nil {
				logrus.WithField("type", "watch").Errorf("Create service: %v", err)
				continue
			}

			if err := s.SetService(srv); err != nil {
				logrus.WithField("type", "watch").Errorf("Set service: %v", err)
				continue
			}

			if !done {
				wg.Done()
				done = true
			}
		}
	}(ch, s)

	wg.Wait()

	return s
}
