package ratelimit

import (
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

var s *gocron.Scheduler = gocron.NewScheduler(time.UTC)

func RunLimitCleaner(l *Limit) error {
	_, err := s.Every(l.Per).Do(func() {
		for k, r := range l.Rates {
			hits := r.Hits
			// Leaky Bucket drop rate
			if hits > 0 && hits > l.MaxRequests {
				Mutex.Lock()
				r.Hits = r.Hits - l.MaxRequests
				Mutex.Unlock()
			} else if hits > 0 && hits <= l.MaxRequests {
				Mutex.Lock()
				r.Hits = 0
				Mutex.Unlock()
			}
			log.Debugf("Run cleaner for key %s, hits %d", k, hits)
		}
	})
	if err != nil {
		log.Errorf("Error adding cleaning job: %s", err)
		return err
	}
	s.StartAsync()
	return nil
}
