package ratelimit

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

var s *gocron.Scheduler = gocron.NewScheduler(time.UTC)

func RunLimitCleaner(l *Limit) {
	_, err := s.Every(l.Per).Do(func() {
		// now := time.Now()
		for k, r := range l.Rates {
			/*
				if r.ExpiredAt.Before(now) {
					l.Rates[k] = createKey()
					fmt.Printf("Run cleaner for key %s\n, expired at %s", k, r.ExpiredAt)
				}
			*/
			//Leaky Bucket drop rate
			if r.Hits > 0 && r.Hits > l.MaxRequests {
				Mutex.Lock()
				r.Hits = r.Hits - l.MaxRequests
				Mutex.Unlock()
			} else if r.Hits > 0 && r.Hits <= l.MaxRequests {
				Mutex.Lock()
				r.Hits = 0
				Mutex.Unlock()
			}
			fmt.Printf("Run cleaner for key %s, hits %d\n", k, r.Hits)
		}

		// fmt.Println("Run cleaner")
	})
	if err != nil {
		fmt.Printf("Error adding cleaning job: %s\n", err)
	}
	s.StartAsync()
	/*
		go func(l *Limit) {
			for {
				time.Sleep(time.Second)
				now := time.Now().Add(-l.Per)
				Mutex.Lock()
				for k, r := range l.Rates {
					if r.ExpiredAt.Before(now) {
						l.Rates[k] = createKey()
					}
				}
				Mutex.Unlock()
			}
		}(l)
	*/
}
