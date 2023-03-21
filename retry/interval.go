package interval

import (
	"time"
)

type RetryFunc func() (bool, error)

type Result struct {
	IsStop            bool
	Err               error
	IsMaxRetryReached bool
}

func Retry(t time.Duration, max int, f RetryFunc) <-chan Result {
	ticker := time.NewTicker(t)
	done := make(chan Result)
	count := 0
	go func() {
		for {
			if count >= max {
				done <- Result{IsStop: true, IsMaxRetryReached: true}
				ticker.Stop()
				close(done)
				return
			}
			select {
			case <-done:
				return
			case <-ticker.C:
				isStop, err := f()
				if err != nil {
					done <- Result{
						IsStop: isStop,
						Err:    err,
					}
					ticker.Stop()
					close(done)
					return
				}
				if isStop {
					ticker.Stop()
					done <- Result{IsStop: isStop}
					close(done)
				}
				count++
			}
		}
	}()
	return done
}
