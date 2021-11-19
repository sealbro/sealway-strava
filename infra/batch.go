package infra

import (
	"sealway-strava/strava"
	"time"
)

func BatchActivities(values <-chan *strava.DetailedActivity, maxItems int, maxTimeout time.Duration) chan []*strava.DetailedActivity {
	batches := make(chan []*strava.DetailedActivity)

	go func() {
		defer close(batches)

		for keepGoing := true; keepGoing; {
			var batch []*strava.DetailedActivity
			expire := time.After(maxTimeout)
			for {
				select {
				case value, ok := <-values:
					if !ok {
						keepGoing = false
						goto done
					}

					batch = append(batch, value)
					if len(batch) == maxItems {
						goto done
					}

				case <-expire:
					goto done
				}
			}

		done:
			if len(batch) > 0 {
				batches <- batch
			}
		}
	}()

	return batches
}
