package ratelimit

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

var s *gocron.Scheduler = gocron.NewScheduler(time.UTC)

func RunLimitCleaner(l *Limit) error {
	_, err := s.Every(l.Per).Do(func() {
		for _, r := range l.Rates {
			Mutex.Lock()
			hits := r.Hits
			// Leaky Bucket drop rate
			if hits > 0 && hits > l.MaxRequests {
				r.Hits = hits - l.MaxRequests
			} else if hits > 0 && hits <= l.MaxRequests {
				r.Hits = 0
			}
			Mutex.Unlock()
		}
	})
	if err != nil {
		panic(fmt.Sprintf("Error adding cleaning job: %s", err))
	}
	s.StartAsync()
	return nil
}
