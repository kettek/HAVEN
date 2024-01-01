package rooms

import "time"

func delayTimeR(d time.Duration) {
	done := make(chan struct{})
	now := time.Now()
	go func() {
		for {
			if time.Since(now) >= d {
				done <- struct{}{}
				return
			}
		}
	}()
	<-done
}
