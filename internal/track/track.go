package track

import (
	"log"
	"time"
)

func Measure(start time.Time, name string, count int) {
	elapsed := time.Since(start)
	log.Printf("Timing: each %s request took %s", name, elapsed/time.Duration(count))
}
