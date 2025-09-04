package track

import (
	"log"
	"time"
)

func Measure(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("Timing: %s took %s", name, elapsed)
}
